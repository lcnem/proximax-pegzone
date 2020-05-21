package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/lcnem/proximax-pegzone/x/proximax-bridge/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/proximax_bridge/unpeg",
		UnpegRequestHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/proximax_bridge/request_invitation",
		UnpegRequestHandlerFn(cliCtx),
	).Methods("POST")
}

type UnpegReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	// TODO: Define more types if needed
	Address              string    `json:"address" yaml:"address"`
	MainchainAddress     string    `json:"mainchain_address" yaml:"mainchain_address"`
	Amount               sdk.Coins `json:"amount" yaml:"amount"`
	FirstCosignerAddress string    `json:"first_cosigner_address" yaml:"first_cosigner_address"`
}

func UnpegRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnpegReq

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, _ := sdk.AccAddressFromBech32(req.Address)
		firstCosignerAddress, _ := sdk.ValAddressFromBech32(req.FirstCosignerAddress)

		// TODO: Define the module tx logic for this action
		msg := types.NewMsgUnpeg(
			address,
			req.MainchainAddress,
			req.Amount,
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

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		address, _ := sdk.ValAddressFromBech32(req.Address)
		firstCosignerAddress, _ := sdk.ValAddressFromBech32(req.FirstCosignerAddress)

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
