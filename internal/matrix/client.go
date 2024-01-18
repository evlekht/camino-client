package matrix

import (
	"bytes"
	"caminoclient/internal/logger"
	"caminoclient/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
)

func NewClient(matrixURI, appServiceAccessToken string, logger logger.Logger, networkID uint32) (*Client, error) {
	return &Client{
		baseURL:               matrixURI + "/_matrix/client/v3",
		appServiceAccessToken: appServiceAccessToken,
		networkID:             networkID,
		hrp:                   constants.GetHRP(networkID),
		logger:                logger,
		utils:                 utils.NewUtils(logger, networkID),
	}, nil
}

type Client struct {
	baseURL               string
	appServiceAccessToken string
	httpClient            http.Client
	logger                logger.Logger
	networkID             uint32
	hrp                   string
	utils                 *utils.UtilsWithLogger
}

type identifier struct {
	Type identifierType `json:"type"`
	User string         `json:"user"`
}

type identifierType string

const identifierTypeUserID identifierType = "m.id.user"

type loginType string

const loginTypeCamino = "m.login.camino"

func (c *Client) Login(key *secp256k1.PrivateKey) (string, string, error) {
	sessionID, payload, err := c.beginLogin(key)
	if err != nil {
		return "", "", err
	}
	return c.login(sessionID, payload, key)
}

func (c *Client) beginLogin(key *secp256k1.PrivateKey) (sessionID string, payload string, err error) {
	keyAddrStr, err := address.Format("t", constants.GetHRP(c.networkID), key.Address().Bytes())
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	request := struct {
		Identifier identifier `json:"identifier"`
	}{
		Identifier: identifier{
			Type: identifierTypeUserID,
			User: keyAddrStr,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	url := fmt.Sprintf("%s/login", c.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		err := fmt.Errorf("beginLogin error: status %d, resp: %s", resp.StatusCode, responseBody)
		c.logger.Error(err)
		return "", "", err
	}

	var response struct {
		Session string `json:"session"`
		Params  struct {
			Camino struct {
				Payload string `json:"payload"`
			} `json:"m.login.camino"`
		} `json:"params"`
		ErrCode string `json:"errcode"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	return response.Session, response.Params.Camino.Payload, nil
}

func (c *Client) login(sessionID, payload string, key *secp256k1.PrivateKey) (string, string, error) {
	signature, err := key.Sign([]byte(payload))
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	encodedSignature, err := formatting.Encode(formatting.Hex, signature)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	keyAddrStr, err := address.Format("t", constants.GetHRP(c.networkID), key.Address().Bytes())
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	type auth struct {
		Session   string    `json:"session"`
		Signature string    `json:"signature"`
		Type      loginType `json:"type"`
	}
	request := struct {
		Identifier identifier `json:"identifier"`
		Auth       auth       `json:"auth"`
	}{
		Identifier: identifier{
			Type: identifierTypeUserID,
			User: keyAddrStr,
		},
		Auth: auth{
			Session:   sessionID,
			Signature: encodedSignature[2:],
			Type:      loginTypeCamino,
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	url := fmt.Sprintf("%s/login", c.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	var response struct {
		UserID      string `json:"user_id"`
		AccessToken string `json:"access_token"`
		HomeServer  string `json:"home_server"`
		DeviceID    string `json:"device_id"`
		ErrCode     string `json:"errcode"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		c.logger.Error(err)
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("login error: status %d, resp: %s", resp.StatusCode, responseBody)
		c.logger.Error(err)
		return "", "", err
	}

	return response.AccessToken, response.DeviceID, nil
}

func (c *Client) CreateRoom() (string, error) {
	request := struct {
		CreationContent struct {
			MFederate bool `json:"m.federate"`
		} `json:"creation_content"`
		Name   string `json:"name"`
		Preset string `json:"preset"`
		Topic  string `json:"topic"`
	}{
		Name:   "The Grand Duke Pub",
		Preset: "public_chat",
		Topic:  "All about happy hour",
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err)
		return "", err
	}

	url := fmt.Sprintf("%s/createRoom", c.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		c.logger.Error(err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+c.appServiceAccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(err)
		return "", err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("create room error: status %d, resp: %s", resp.StatusCode, responseBody)
		c.logger.Error(err)
		return "", err
	}

	var response struct {
		RoomID string `json:"room_id"`
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		c.logger.Error(err)
		return "", err
	}

	c.logger.Debugf("Created room: %s", response.RoomID)

	return response.RoomID, nil
}

// func (c *Client) SendMessageWithCheque(roomID, message string, cheque *tTxs.SignedCheque) error {
// 	signatures := cheque.Auth.Signatures()
// 	if len(signatures) != 1 {
// 		err := fmt.Errorf("unexpected number of signatures: expected 1, but got %d", len(signatures))
// 		c.logger.Error(err)
// 		return err
// 	}

// 	// encodedSignature, err := formatting.Encode(formatting.Hex, signatures[0][:])
// 	// if err != nil {
// 	// 	c.logger.Error(err)
// 	// 	return err
// 	// }

// 	issuer, err := address.Format("T", c.hrp, cheque.Issuer.Bytes())
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	beneficiary, err := address.Format("T", c.hrp, cheque.Beneficiary.Bytes())
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	type msgCheque struct {
// 		Issuer      string            `json:"issuer"`
// 		Agent       string            `json:"agent"`
// 		Beneficiary string            `json:"beneficiary"`
// 		Amount      json_utils.Uint64 `json:"amount"`
// 		SerialID    json_utils.Uint64 `json:"serialID"`
// 		Signature   string            `json:"signature"`
// 	}

// 	request := struct {
// 		Body    string      `json:"body"`
// 		MsgType string      `json:"msgtype"`
// 		Cheques []msgCheque `json:"cheques"`
// 	}{
// 		Body:    message,
// 		MsgType: "C4TContentHotelAvailRequest",
// 		Cheques: []msgCheque{{
// 			Issuer:      issuer,
// 			Agent:       cheque.Agent.String(),
// 			Beneficiary: beneficiary,
// 			Amount:      json_utils.Uint64(cheque.Amount),
// 			SerialID:    json_utils.Uint64(cheque.SerialID),
// 			Signature:   "e005367ebfaee8a9d980995891c332ff31bc7fb17dc9727c927dec847486bea80a59b22b398c454fb1505b34633a7faed16f1fe9e96a14157332ec6f8c30cb6701",
// 			// Signature:   encodedSignature[2:],
// 		}},
// 	}

// 	requestBody, err := json.Marshal(request)
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	uuid, err := uuid.NewRandom()
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	url := fmt.Sprintf("%s/rooms/%s/send/m.room.c4t-msg/%s", c.baseURL, roomID, uuid.String())
// 	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}
// 	req.Header.Add("Authorization", "Bearer "+c.appServiceAccessToken)

// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	responseBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		c.logger.Error(err)
// 		return err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		err := fmt.Errorf("SendMessageWithCheque error: status %d, resp: %s", resp.StatusCode, responseBody)
// 		c.logger.Error(err)
// 		return err
// 	}

// 	return nil
// }
