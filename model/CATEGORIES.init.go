package model

func init() {
	// init first level
	CategoryAudio.Children = Categories{
		CategoryAudioAmplifiersMixers,
		CategoryAudioAudioVideoCablesConverters,
		CategoryAudioEarphonesHeadphonesHeadsets,
		CategoryAudioHomeAudioSpeakers,
		CategoryAudioMediaPlayers,
		CategoryAudioMicrophones,
		CategoryAudioOther,
	}
	CategoryAutomobiles.Children = Categories{
		CategoryAutomobilesAutomobileExteriorAccessories,
		CategoryAutomobilesAutomobileInteriorAccessories,
		CategoryAutomobilesAutomobileSpareParts,
		CategoryAutomobilesAutomotiveCare,
		CategoryAutomobilesAutomotiveKeychainsKeyCovers,
		CategoryAutomobilesAutomotiveOilsLubes,
		CategoryAutomobilesAutomotiveTools,
	}
	CategoryBabyKidsFashion.Children = Categories{
		CategoryBabyKidsFashionBabyClothes,
		CategoryBabyKidsFashionBabyKidsAccessories,
		CategoryBabyKidsFashionBabyMittensFootwear,
		CategoryBabyKidsFashionBoyClothes,
		CategoryBabyKidsFashionBoyShoes,
		CategoryBabyKidsFashionGirlClothes,
		CategoryBabyKidsFashionGirlShoes,
		CategoryBabyKidsFashionUnderwearInnerwear,
	}
	CategoryBeauty.Children = Categories{
		CategoryBeautyBathBodyCare,
		CategoryBeautyBeautySetsPackages,
		CategoryBeautyBeautyTools,
		CategoryBeautyHairCare,
		CategoryBeautyHandFootNailCare,
		CategoryBeautyMakeup,
		CategoryBeautyMensCare,
		CategoryBeautyPerfumesFragrances,
		CategoryBeautySkincare,
	}
	CategoryBooksMagazines.Children = Categories{
		CategoryBooksMagazinesBooks,
		CategoryBooksMagazinesEBooks,
		CategoryBooksMagazinesMagazinesNewspaper,
	}
	CategoryCamerasDrones.Children = Categories{
		CategoryCamerasDronesCameraAccessories,
		CategoryCamerasDronesCameraCare,
		CategoryCamerasDronesCameras,
		CategoryCamerasDronesDroneAccessories,
		CategoryCamerasDronesDrones,
		CategoryCamerasDronesLensAccessories,
		CategoryCamerasDronesLenses,
		CategoryCamerasDronesSecurityCamerasSystems,
	}
	CategoryComputersAccessories.Children = Categories{
		CategoryComputersAccessoriesDataStorage,
		CategoryComputersAccessoriesDesktopComputer,
		CategoryComputersAccessoriesDesktopLaptopComponents,
		CategoryComputersAccessoriesKeyboardsMice,
		CategoryComputersAccessoriesLaptop,
		CategoryComputersAccessoriesMonitors,
		CategoryComputersAccessoriesNetworkComponents,
		CategoryComputersAccessoriesOfficeEquipment,
		CategoryComputersAccessoriesPeripheralsAccessories,
		CategoryComputersAccessoriesPrintersScanners,
		CategoryComputersAccessoriesSoftwares,
	}
	CategoryFashionAccessories.Children = Categories{
		CategoryFashionAccessoriesAccessoriesSetsPackages,
		CategoryFashionAccessoriesAdditionalAccessories,
		CategoryFashionAccessoriesAnklets,
		CategoryFashionAccessoriesBelts,
		CategoryFashionAccessoriesBraceletsBangles,
		CategoryFashionAccessoriesEarrings,
		CategoryFashionAccessoriesEyewear,
		CategoryFashionAccessoriesGloves,
		CategoryFashionAccessoriesHairAccessories,
		CategoryFashionAccessoriesHatsCaps,
		CategoryFashionAccessoriesInvestmentPreciousMetals,
		CategoryFashionAccessoriesNecklaces,
		CategoryFashionAccessoriesNecktiesBowTiesCravats,
		CategoryFashionAccessoriesRings,
		CategoryFashionAccessoriesScarvesShawls,
	}
	CategoryFoodBeverage.Children = Categories{
		CategoryFoodBeverageBakingNeeds,
		CategoryFoodBeverageBeverages,
		CategoryFoodBeverageBreakfastCerealsSpread,
		CategoryFoodBeverageConvenienceReadytoeat,
		CategoryFoodBeverageCookingEssentials,
		CategoryFoodBeverageDairyEggs,
		CategoryFoodBeverageFoodStaples,
		CategoryFoodBeverageSeasoningsCondiments,
		CategoryFoodBeverageSnacks,
	}
	CategoryGamingConsoles.Children = Categories{
		CategoryGamingConsolesConsoleAccessories,
		CategoryGamingConsolesConsoleMachines,
		CategoryGamingConsolesVideoGames,
	}
	CategoryHealth.Children = Categories{
		CategoryHealthFoodSupplement,
		CategoryHealthMedicalSupplies,
		CategoryHealthPersonalCare,
		CategoryHealthSexualWellness,
	}
	CategoryHobbiesCollections.Children = Categories{
		CategoryHobbiesCollectionsCDDVDBluray,
		CategoryHobbiesCollectionsCollectibleItems,
		CategoryHobbiesCollectionsMusicalInstrumentsAccessories,
		CategoryHobbiesCollectionsNeedlework,
		CategoryHobbiesCollectionsPhotoAlbums,
		CategoryHobbiesCollectionsSouvenirs,
		CategoryHobbiesCollectionsToysGames,
		CategoryHobbiesCollectionsVinylRecords,
	}
	CategoryHomeAppliances.Children = Categories{
		CategoryHomeAppliancesBatteries,
		CategoryHomeAppliancesElectricalCircuitryParts,
		CategoryHomeAppliancesKitchenAppliances,
		CategoryHomeAppliancesLargeHouseholdAppliances,
		CategoryHomeAppliancesProjectorsAccessories,
		CategoryHomeAppliancesRemoteControls,
		CategoryHomeAppliancesSmallHouseholdAppliances,
		CategoryHomeAppliancesTVsAccessories,
	}
	CategoryHomeLiving.Children = Categories{
		CategoryHomeLivingBathrooms,
		CategoryHomeLivingBedding,
		CategoryHomeLivingDecoration,
		CategoryHomeLivingDinnerware,
		CategoryHomeLivingFengshuiReligiousSupplies,
		CategoryHomeLivingFurniture,
		CategoryHomeLivingGardening,
		CategoryHomeLivingHandWarmersHotWaterBagsIceBags,
		CategoryHomeLivingHomeCareSupplies,
		CategoryHomeLivingHomeFragranceAromatherapy,
		CategoryHomeLivingHomeOrganizers,
		CategoryHomeLivingKitchenware,
		CategoryHomeLivingLighting,
		CategoryHomeLivingPartySupplies,
		CategoryHomeLivingSafetySecurity,
		CategoryHomeLivingToolsHomeImprovement,
	}
	CategoryMenBags.Children = Categories{
		CategoryMenBagsBackpacks,
		CategoryMenBagsBriefcases,
		CategoryMenBagsClutches,
		CategoryMenBagsCrossbodyShoulderBags,
		CategoryMenBagsLaptopBags,
		CategoryMenBagsToteBags,
		CategoryMenBagsWaistBagsChestBags,
		CategoryMenBagsWallets,
	}
	CategoryMenClothes.Children = Categories{
		CategoryMenClothesCostumes,
		CategoryMenClothesHoodiesSweatshirts,
		CategoryMenClothesInnerwearUnderwear,
		CategoryMenClothesJacketsCoatsVests,
		CategoryMenClothesJeans,
		CategoryMenClothesOccupationalAttire,
		CategoryMenClothesPants,
		CategoryMenClothesSets,
		CategoryMenClothesShorts,
		CategoryMenClothesSleepwear,
		CategoryMenClothesSocks,
		CategoryMenClothesSuits,
		CategoryMenClothesSweatersCardigans,
		CategoryMenClothesTops,
		CategoryMenClothesTraditionalWear,
		CategoryMenClothesWinterJacketsCoats,
	}
	CategoryMenShoes.Children = Categories{
		CategoryMenShoesBoots,
		CategoryMenShoesLoafersBoatShoes,
		CategoryMenShoesOxfordsLaceUps,
		CategoryMenShoesSandalsFlipFlops,
		CategoryMenShoesShoeCareAccessories,
		CategoryMenShoesSlipOnsMules,
		CategoryMenShoesSneakers,
	}
	CategoryMobileGadgets.Children = Categories{
		CategoryMobileGadgetsAccessories,
		CategoryMobileGadgetsMobilePhones,
		CategoryMobileGadgetsOther,
		CategoryMobileGadgetsSimCards,
		CategoryMobileGadgetsTablets,
		CategoryMobileGadgetsWalkieTalkies,
		CategoryMobileGadgetsWearableDevices,
	}
	CategoryMomBaby.Children = Categories{
		CategoryMomBabyBabyHealthcare,
		CategoryMomBabyBabySafety,
		CategoryMomBabyBabyTravelEssentials,
		CategoryMomBabyBathBodyCare,
		CategoryMomBabyDiaperingPotty,
		CategoryMomBabyFeedingEssentials,
		CategoryMomBabyGiftSetsPackages,
		CategoryMomBabyMaternityAccessories,
		CategoryMomBabyMaternityHealthcare,
		CategoryMomBabyMilkFormulaBabyFood,
		CategoryMomBabyNursery,
		CategoryMomBabyToys,
	}
	CategoryMotorcycles.Children = Categories{
		CategoryMotorcyclesMotorcycleAccessories,
		CategoryMotorcyclesMotorcycleHelmetsAccessories,
		CategoryMotorcyclesMotorcycleSpareParts,
		CategoryMotorcyclesMotorcycles,
		CategoryMotorcyclesOther,
	}
	CategoryPets.Children = Categories{
		CategoryPetsLitterToilet,
		CategoryPetsPetAccessories,
		CategoryPetsPetClothingAccessories,
		CategoryPetsPetFood,
		CategoryPetsPetGrooming,
		CategoryPetsPetHealthcare,
	}
	CategorySportsOutdoors.Children = Categories{
		CategorySportsOutdoorsSportFootwear,
		CategorySportsOutdoorsSportsOutdoorAccessories,
		CategorySportsOutdoorsSportsOutdoorApparels,
		CategorySportsOutdoorsSportsOutdoorRecreationEquipments,
	}
	CategoryStationery.Children = Categories{
		CategoryStationeryArtSupplies,
		CategoryStationeryGiftWrapping,
		CategoryStationeryLettersEnvelopes,
		CategoryStationeryNotebooksPapers,
		CategoryStationeryOther,
		CategoryStationerySchoolOfficeEquipment,
		CategoryStationeryWritingCorrection,
	}
	CategoryTravelLuggage.Children = Categories{
		CategoryTravelLuggageLuggage,
		CategoryTravelLuggageTravelAccessories,
		CategoryTravelLuggageTravelBags,
	}
	CategoryWatches.Children = Categories{
		CategoryWatchesMenWatches,
		CategoryWatchesSetCoupleWatches,
		CategoryWatchesWatchesAccessories,
		CategoryWatchesWomenWatches,
	}
	CategoryWomenBags.Children = Categories{
		CategoryWomenBagsBackpacks,
		CategoryWomenBagsBagAccessories,
		CategoryWomenBagsClutchesWristlets,
		CategoryWomenBagsCrossbodyShoulderBags,
		CategoryWomenBagsLaptopBags,
		CategoryWomenBagsTophandleBags,
		CategoryWomenBagsToteBags,
		CategoryWomenBagsWaistBagsChestBags,
		CategoryWomenBagsWallets,
	}
	CategoryWomenClothes.Children = Categories{
		CategoryWomenClothesCostumes,
		CategoryWomenClothesDresses,
		CategoryWomenClothesFabric,
		CategoryWomenClothesHoodiesSweatshirts,
		CategoryWomenClothesJacketsCoatsVests,
		CategoryWomenClothesJeans,
		CategoryWomenClothesJumpsuitsPlaysuitsOveralls,
		CategoryWomenClothesLingerieUnderwear,
		CategoryWomenClothesMaternityWear,
		CategoryWomenClothesPantsLeggings,
		CategoryWomenClothesSets,
		CategoryWomenClothesShorts,
		CategoryWomenClothesSkirts,
		CategoryWomenClothesSleepwearPajamas,
		CategoryWomenClothesSocksStockings,
		CategoryWomenClothesSweatersCardigans,
		CategoryWomenClothesTops,
		CategoryWomenClothesTraditionalWear,
		CategoryWomenClothesWeddingDresses,
	}
	CategoryWomenShoes.Children = Categories{
		CategoryWomenShoesBoots,
		CategoryWomenShoesFlatSandalsFlipFlops,
		CategoryWomenShoesFlats,
		CategoryWomenShoesHeels,
		CategoryWomenShoesShoeCareAccessories,
		CategoryWomenShoesSneakers,
		CategoryWomenShoesWedges,
	}
}
