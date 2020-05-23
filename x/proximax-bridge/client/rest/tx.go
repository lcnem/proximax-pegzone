package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/proximax_bridge/peg",
		PegRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/proximax_bridge/unpeg",
		UnpegRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/proximax_bridge/request_invitation",
		UnpegRequestHandlerFn(cliCtx),
	).Methods("POST")
}

type PegReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
	Address         string `json:"address" yaml:"address"`
	MainchainTxHash string `json:"mainchain_tx_hash" yaml:"mainchain_tx_hash"`
	Amount          string `json:"amount" yaml:"amount"`
}

func PegRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PegReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			msg := fmt.Sprintf("failed to parse address: %s", req.Address)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		if len(strings.Trim(req.MainchainTxHash, "")) == 0 {
			msg := fmt.Sprintf("invalid mainchain_tx_hash: %s", req.MainchainTxHash)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			msg := fmt.Sprintf("failed to parse amount: %s", req.Amount)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		msg := types.NewMsgPeg(address, req.MainchainTxHash, amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type UnpegReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
	Address              string `json:"address" yaml:"address"`
	MainchainAddress     string `json:"mainchain_address" yaml:"mainchain_address"`
	Amount               string `json:"amount" yaml:"amount"`
	FirstCosignerAddress string `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

func UnpegRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnpegReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.AccAddressFromBech32(req.Address)
		if err != nil {
			msg := fmt.Sprintf("failed to parse address: %s", req.Address)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		if len(strings.Trim(req.MainchainAddress, "")) == 0 {
			msg := fmt.Sprintf("invalid mainchain_address: %s", req.MainchainAddress)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			msg := fmt.Sprintf("failed to parse amount: %s", req.Amount)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		firstCosignerAddress, err := sdk.ValAddressFromBech32(req.FirstCosignerAddress)
		if err != nil {
			msg := fmt.Sprintf("failed to parse first_cosigner_address: %s", req.FirstCosignerAddress)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		// TODO: Define the module tx logic for this action
		msg := types.NewMsgUnpeg(
			address,
			req.MainchainAddress,
			amount,
			firstCosignerAddress,
		)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type RequestInvitationReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
	Address              string `json:"address" yaml:"address"`
	NewCosignerPublicKey string `json:"new_cosigner_public_key" yaml:"new_cosigner_public_key"`
	FirstCosignerAddress string `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

func RequestInvitationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestInvitationReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, err := sdk.ValAddressFromBech32(req.Address)
		if err != nil {
			msg := fmt.Sprintf("failed to parse address: %s", req.Address)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		if len(strings.Trim(req.NewCosignerPublicKey, "")) == 0 {
			msg := fmt.Sprintf("invalid new_cosigner_public_key: %s", req.NewCosignerPublicKey)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		firstCosignerAddress, err := sdk.ValAddressFromBech32(req.FirstCosignerAddress)
		if err != nil {
			msg := fmt.Sprintf("failed to parse first_cosigner_address: %s", req.FirstCosignerAddress)
			rest.WriteErrorResponse(w, http.StatusBadRequest, msg)
			return
		}

		// TODO: Define the module tx logic for this action
		msg := types.NewMsgRequestInvitation(
			address,
			req.NewCosignerPublicKey,
			firstCosignerAddress,
		)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

/*
// Action TX body
type <Action>Req struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
}

func <Action>RequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req <Action>Req
		vars := mux.Vars(r)

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// TODO: Define the module tx logic for this action

		utils.WriteGenerateStdTxResponse(w, cliCtx, BaseReq, []sdk.Msg{msg})
	}
}
*/
