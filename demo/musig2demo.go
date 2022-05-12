package main

import (
	"log"

	"github.com/chainx-org/bitcoin-go-api/musig2"
)

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.SetPrefix("[Bitcoin-Taproot]")
	// Generate non-threshold signature address
	PHRASE0 := "flame flock chunk trim modify raise rough client coin busy income smile"
	private0, err := musig2.GetMyPrivkey(PHRASE0, "")
	if err != nil {
		log.Fatal(err)
	}
	pubkey0, err := musig2.GetMyPubkey(private0)
	if err != nil {
		log.Fatal(err)
	}
	address0, err := musig2.GetMyAddress(pubkey0, "signet")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("address0: ", address0)

	PHRASE1 := "shrug argue supply evolve alarm caught swamp tissue hollow apology youth ethics"
	private1, err := musig2.GetMyPrivkey(PHRASE1, "")
	if err != nil {
		log.Fatal(err)
	}
	pubkey1, err := musig2.GetMyPubkey(private1)
	if err != nil {
		log.Fatal(err)
	}

	PHRASE2 := "awesome beef hill broccoli strike poem rebel unique turn circle cool system"
	private2, err := musig2.GetMyPrivkey(PHRASE2, "")
	if err != nil {
		log.Fatal(err)
	}
	pubkey2, err := musig2.GetMyPubkey(private2)
	if err != nil {
		log.Fatal(err)
	}

	// Cost of non-threshold signature addresses
	prevTxs := []string{"020000000001014be640313b023c3c731b7e89c3f97bebcebf9772ea2f7747e5604f4483a447b601000000000000000002a0860100000000002251209a9ea267884f5549c206b2aec2bd56d98730f90532ea7f7154d4d4f923b7e3bbc027090000000000225120c9929543dfa1e0bb84891acd47bfa6546b05e26b7a04af8eb6765fcc969d565f01404dc68b31efc1468f84db7e9716a84c19bbc53c2d252fd1d72fa6469e860a74486b0990332b69718dbcb5acad9d48634d23ee9c215ab15fb16f4732bed1770fdf00000000"}
	txids := []string{"1f8e0f7dfa37b184244d022cdf2bc7b8e0bac8b52143ea786fa3f7bbe049eeae"}
	inputIndexs := []uint32{1}
	addresses := []string{"tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw", "35516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38", "tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68"}
	amounts := []uint64{100000, 0, 400000}
	baseTx, err := musig2.GenerateRawTx(prevTxs, txids, inputIndexs, addresses, amounts)
	if err != nil {
		log.Fatal(err)
	}
	unsignedTx, err := musig2.GetUnsignedTx(baseTx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Unsigned Tx: ", unsignedTx)
	var taprootTx string
	for i := 0; i < len(txids); i++ {
		privkey := "4a84a4601e463bc02dd0b8be03f3721187e9fc3105d5d5e8930ff3c8ca15cf40"
		sighash, err := musig2.GetSighash(baseTx, txids[i], inputIndexs[i], "", 0)
		log.Println("Sighash: ", sighash)
		if err != nil {
			log.Fatal(err)
		}
		schnorrSignature, err := musig2.GenerateSchnorrSignature(sighash, privkey)
		log.Println("SchnorrSignature: ", schnorrSignature)
		if err != nil {
			log.Fatal(err)
		}
		taprootTx, err = musig2.BuildTaprootTx(baseTx, schnorrSignature, txids[i], inputIndexs[i])
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Current Taproot Tx: ", taprootTx)
	}
	log.Println("Final Taproot Tx: ", taprootTx)

	// Generate threshold signature address
	thresholdPubkey, err := musig2.GenerateThresholdPubkey([]string{pubkey0, pubkey1, pubkey2}, 2)
	if err != nil {
		log.Fatal(err)
	}
	thresholdAddress, err := musig2.GetMyAddress(thresholdPubkey, "signet")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Threshold Address: ", thresholdAddress)

	// Cost of threshold signature address
	privateA := "e5bb018d70c6fb5dd8ad91f6c88fb0e6fdab2c482978c95bb3794ca6e2e50dc2"
	privateB := "a7150e8f24ab26ebebddd831aeb8f00ecb593df3b80ae1e8b8be01351805f2d6"
	privateC := "4a84a4601e463bc02dd0b8be03f3721187e9fc3105d5d5e8930ff3c8ca15cf40"
	pubkeyA, err := musig2.GetMyPubkey(privateA)
	if err != nil {
		log.Fatal(err)
	}
	pubkeyB, err := musig2.GetMyPubkey(privateB)
	if err != nil {
		log.Fatal(err)
	}
	pubkeyC, err := musig2.GetMyPubkey(privateC)
	if err != nil {
		log.Fatal(err)
	}

	prevTxs = []string{"02000000000101aeee49e0bbf7a36f78ea4321b5c8bae0b8c72bdf2c024d2484b137fa7d0f8e1f01000000000000000003a0860100000000002251209a9ea267884f5549c206b2aec2bd56d98730f90532ea7f7154d4d4f923b7e3bb0000000000000000326a3035516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a38801a060000000000225120c9929543dfa1e0bb84891acd47bfa6546b05e26b7a04af8eb6765fcc969d565f01409e325889515ed47099fdd7098e6fafdc880b21456d3f368457de923f4229286e34cef68816348a0581ae5885ede248a35ac4b09da61a7b9b90f34c200872d2e300000000"}
	txids = []string{"8e5d37c768acc4f3e794a10ad27bf0256237c80c22fa67117e3e3e1aec22ea5f"}
	inputIndexs = []uint32{0}
	addresses = []string{"tb1pexff2s7l58sthpyfrtx500ax234stcnt0gz2lr4kwe0ue95a2e0srxsc68", "tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw"}
	amounts = []uint64{50000, 40000}
	baseTx, err = musig2.GenerateRawTx(prevTxs, txids, inputIndexs, addresses, amounts)
	if err != nil {
		log.Fatal(err)
	}
	unsignedTx, err = musig2.GetUnsignedTx(baseTx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Unsigned Tx: ", unsignedTx)

	var thresholdTx string
	for i := 0; i < len(txids); i++ {
		pubkeyBC, err := musig2.GetAggPublicKey([]string{pubkeyB, pubkeyC})
		if err != nil {
			log.Fatal(err)
		}
		sighash, err := musig2.GetSighash(baseTx, txids[i], inputIndexs[i], pubkeyBC, 1)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Sighash: ", sighash)
		round1State0 := musig2.GetRound1State()
		round1State0Encode, err := musig2.EncodeRound1State(round1State0)
		if err != nil {
			log.Fatal(err)
		}
		round1State0 = musig2.DecodeRound1State(round1State0Encode)
		round1State1 := musig2.GetRound1State()

		round1Msg0, err := musig2.GetRound1Msg(round1State0)
		if err != nil {
			log.Fatal(err)
		}
		round1Msg1, err := musig2.GetRound1Msg(round1State1)
		if err != nil {
			log.Fatal(err)
		}
		round2Msg0, err := musig2.GetRound2Msg(round1State0, sighash, privateB, []string{pubkeyB, pubkeyC}, []string{round1Msg1})
		if err != nil {
			log.Fatal(err)
		}
		round2Msg1, err := musig2.GetRound2Msg(round1State1, sighash, privateC, []string{pubkeyB, pubkeyC}, []string{round1Msg0})
		if err != nil {
			log.Fatal(err)
		}
		multiSignature, err := musig2.GetAggSignature([]string{round2Msg0, round2Msg1})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("MultiSignature: ", multiSignature)
		controlBlock, err := musig2.GenerateControlBlock([]string{pubkeyA, pubkeyB, pubkeyC}, 2, pubkeyBC)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Control Block: ", controlBlock)
		thresholdTx, err = musig2.BuildThresholdTx(baseTx, multiSignature, pubkeyBC, controlBlock, txids[i], inputIndexs[i])
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Current Threshold Tx: ", thresholdTx)
	}
	log.Println("Final Threshold Tx: ", thresholdTx)

	// other tool func test
	scriptPubkey, err := musig2.GetScriptPubkey("tb1pn202yeugfa25nssxk2hv902kmxrnp7g9xt487u256n20jgahuwasdcjfdw")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ScriptPubkey: ", scriptPubkey)
}
