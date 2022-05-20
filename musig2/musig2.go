package musig2

/*
#cgo darwin LDFLAGS: -L../lib/mac -lmusig2_dll
#cgo linux LDFLAGS: -L../lib/linux -lmusig2_dll
#cgo windows LDFLAGS: -L../lib/windows -lmusig2_dll
#include <stdlib.h>
#include "../lib/Musig2Header.h"
*/
import "C"

import (
	"encoding/hex"
	"errors"
	"strings"
	"unsafe"
)

// Just check if the result is hex string
func verifyResult(result *C.char) (string, error) {
	output := C.GoString(result)
	_, err := hex.DecodeString(output)
	if err != nil {
		return "", errors.New(output)
	} else {
		return output, nil
	}
}

// Generate the private key from the mnemonic phrases,
// and the default derived password is empty.
func GetMyPrivkey(phrase string, passphrase string) (string, error) {
	cPhrase := C.CString(phrase)
	defer C.free(unsafe.Pointer(cPhrase))
	cPassphrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cPassphrase))
	result := C.get_my_privkey_musig2(cPhrase, cPassphrase)
	return verifyResult(result)
}

// Generate the corresponding public key from the private key.
func GetMyPubkey(priv string) (string, error) {
	cPriv := C.CString(priv)
	defer C.free(unsafe.Pointer(cPriv))
	result := C.get_my_pubkey_musig2(cPriv)
	return verifyResult(result)
}

// No parameters are needed to get the state pointer for the
// first round. Note that manual release may be required at
// the upper level.
func GetRound1State() *C.State {
	return C.get_round1_state()
}

// Generate bitcoin addresses
//
// network supports "mainnet", "testnet", "regtest", "signet"
func GetMyAddress(pubkey string, network string) (string, error) {
	cPubkey := C.CString(pubkey)
	defer C.free(unsafe.Pointer(cPubkey))
	cNetwork := C.CString(network)
	defer C.free(unsafe.Pointer(cNetwork))
	result := C.generate_btc_address(cPubkey, cNetwork)
	return C.GoString(result), nil
}

// Get the first round of messages to be broadcast
func GetRound1Msg(state *C.State) (string, error) {
	result := C.get_round1_msg(state)
	return verifyResult(result)
}

// Encode the first round of state for easy local persistent storage
func EncodeRound1State(state *C.State) (string, error) {
	result := C.encode_round1_state(state)
	return C.GoString(result), nil
}

// Parasing the first round of state from the local persistent store
func DecodeRound1State(state string) *C.State {
	cState := C.CString(state)
	defer C.free(unsafe.Pointer(cState))
	result := C.decode_round1_state(cState)
	return result
}

// Get the second round of messages to be broadcast
func GetRound2Msg(state *C.State, msg string, priv string, pubkeys []string, receivedRound1Msg []string) (string, error) {
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))
	cPriv := C.CString(priv)
	defer C.free(unsafe.Pointer(cPriv))
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	allRound1Msgs := strings.Join(receivedRound1Msg, "")
	cAllRound1Msgs := C.CString(allRound1Msgs)
	defer C.free(unsafe.Pointer(cAllRound1Msgs))
	result := C.get_round2_msg(state, cMsg, cPriv, cAllPubkeys, cAllRound1Msgs)
	return verifyResult(result)
}

// Generate the final signature using all the messages from the second round
func GetAggSignature(round2Msg []string) (string, error) {
	allRound2Msg := strings.Join(round2Msg, "")
	cAllRound2Msg := C.CString(allRound2Msg)
	defer C.free(unsafe.Pointer(cAllRound2Msg))
	result := C.get_signature_musig2(cAllRound2Msg)
	return verifyResult(result)
}

// Aggregate multiple public keys into one public key
//
// pubkey is full public key, 65 bytes
func GetAggPublicKey(pubkeys []string) (string, error) {
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	result := C.get_key_agg(cAllPubkeys)
	return verifyResult(result)
}

// Generate threshold signature public key
func GenerateThresholdPubkey(pubkeys []string, threshold uint8) (string, error) {
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	result := C.generate_threshold_pubkey_musig2(cAllPubkeys, C.uint8_t(threshold))
	return verifyResult(result)
}

// Generate a proof of the aggregated public key by
// passing in the public key and signature threshold of
// all signers and the aggregated public key of everyone
// who performed the signature this time.
func GenerateControlBlock(pubkeys []string, threshold uint8, aggPubkey string) (string, error) {
	allPubkeys := strings.Join(pubkeys, "")
	cAllPubkeys := C.CString(allPubkeys)
	defer C.free(unsafe.Pointer(cAllPubkeys))
	cAggPubkey := C.CString(aggPubkey)
	defer C.free(unsafe.Pointer(cAggPubkey))
	result := C.generate_control_block_musig2(cAllPubkeys, C.uint8_t(threshold), cAggPubkey)
	return verifyResult(result)
}

// Generate schnorr signature.
func GenerateSchnorrSignature(message string, privkey string) (string, error) {
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	cPrivkey := C.CString(privkey)
	defer C.free(unsafe.Pointer(cPrivkey))
	result := C.generate_schnorr_signature(cMessage, cPrivkey)
	return verifyResult(result)
}

// Generate script pubkey from address
func GetScriptPubkey(addr string) (string, error) {
	cAddr := C.CString(addr)
	defer C.free(unsafe.Pointer(cAddr))
	result := C.get_script_pubkey(cAddr)
	return verifyResult(result)
}

// Add the first input[txid + outpoint's index] to initialize basic transactions
func GenerateRawTx(prevTxs []string, txids []string, inputIndexs []uint32, addresses []string, amounts []uint64) (string, error) {
	if len(txids) != len(inputIndexs) {
		return "", errors.New("txids and inputIndexs must be the same length")
	}

	if len(txids) != len(prevTxs) {
		return "", errors.New("txids and prevTxs must be the same length")
	}

	if len(addresses) != len(amounts) {
		return "", errors.New("addresses and amounts must be the same length")
	}

	if len(txids) == 0 {
		return "", errors.New("must provide at least one txid")
	}

	if len(addresses) == 0 {
		return "", errors.New("must provide at least one address")
	}

	cPrevTx := C.CString(prevTxs[0])
	defer C.free(unsafe.Pointer(cPrevTx))
	cTxid := C.CString(txids[0])
	defer C.free(unsafe.Pointer(cTxid))

	baseTx := C.get_base_tx(cPrevTx, cTxid, C.uint32_t(inputIndexs[0]))
	for i := 1; i < len(txids); i++ {
		cPrevTx = C.CString(prevTxs[i])
		cTxid = C.CString(txids[i])
		baseTx = C.add_input(baseTx, cPrevTx, cTxid, C.uint32_t(inputIndexs[i]))
	}

	var cAddress *C.char
	defer C.free(unsafe.Pointer(cAddress))
	for i := 0; i < len(addresses); i++ {
		cAddress = C.CString(addresses[i])
		baseTx = C.add_output(baseTx, cAddress, C.uint64_t(amounts[i]))
	}
	return verifyResult(baseTx)
}

// Calculate the sighash of the transaction input
//
// Passing the previous transaction and the constructed transaction
// and the previous transaction outpoint index to calculate Sighash.
//
// [`agg_pubkey`]: required when spending by script, pass in a null value when spending by path
// [`sigversion`]: 0 or 1, 0 is Taproot, 1 is Tapscript.
func GetSighash(tx string, txid string, inputIndex uint32, aggPubkey string, sigversion uint32) (string, error) {
	cTx := C.CString(tx)
	defer C.free(unsafe.Pointer(cTx))
	cTxid := C.CString(txid)
	defer C.free(unsafe.Pointer(cTxid))
	cAggPubkey := C.CString(aggPubkey)
	defer C.free(unsafe.Pointer(cAggPubkey))
	result := C.get_sighash(cTx, cTxid, C.uint32_t(inputIndex), cAggPubkey, C.uint32_t(sigversion))
	return verifyResult(result)
}

// Construct Threshold address spending transaction.
//
// [`base_tx`]: tx with at least one input and one output.
// [`agg_signature`]: aggregate signature of sighash
// [`agg_pubkey`]: signature corresponding to the aggregate public key.
// [`control`]: control script.
// [`input_index`]: index of the input in base_tx.
func BuildThresholdTx(tx string, aggSignature string, aggPubkey string, control string, txid string, inputIndex uint32) (string, error) {
	cTx := C.CString(tx)
	defer C.free(unsafe.Pointer(cTx))
	cAggSignature := C.CString(aggSignature)
	defer C.free(unsafe.Pointer(cAggSignature))
	cAggPubkey := C.CString(aggPubkey)
	defer C.free(unsafe.Pointer(cAggPubkey))
	cControl := C.CString(control)
	defer C.free(unsafe.Pointer(cControl))
	cTxid := C.CString(txid)
	defer C.free(unsafe.Pointer(cTxid))
	result := C.build_raw_script_tx(cTx, cAggSignature, cAggPubkey, cControl, cTxid, C.uint32_t(inputIndex))
	return verifyResult(result)
}

// Construct normal taproot address spending transaction.
//
// [`base_tx`]: tx with at least one input and one output.
// [`signature`]: signature of sighash
// [`input_index`]: index of the input in base_tx.
func BuildTaprootTx(tx string, signature string, txid string, inputIndex uint32) (string, error) {
	cTx := C.CString(tx)
	defer C.free(unsafe.Pointer(cTx))
	cSignature := C.CString(signature)
	defer C.free(unsafe.Pointer(cSignature))
	cTxid := C.CString(txid)
	defer C.free(unsafe.Pointer(cTxid))
	result := C.build_raw_key_tx(cTx, cSignature, cTxid, C.uint32_t(inputIndex))
	return verifyResult(result)
}

// Obtain the original unsigned transaction for testing the correctness of the transaction
func GetUnsignedTx(tx string) (string, error) {
	cTx := C.CString(tx)
	defer C.free(unsafe.Pointer(cTx))
	result := C.get_unsigned_tx(cTx)
	return verifyResult(result)
}
