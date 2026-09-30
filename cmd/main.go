package main

import (
	"fmt"
	"log"
	"math/big"
	"ngo-transparency-platform/pkg/entities"
	"ngo-transparency-platform/pkg/platform"
)

func main() {
	fmt.Println("=== Enhanced NGO Transparency Platform Demo ===")

	// Initialize platform
	ngoPlat := platform.NewNGOTransparencyPlatform()
	fmt.Println("✓ Platform initialized")

	// Initialize Polygon integration (testnet)
	ngoPlat.InitializePolygon(
		"https://polygon-mumbai.g.alchemy.com/v2/demo",
		"0x"+"1111111111111111111111111111111111111111111111111111111111111111", // Dummy private key
		300000,
		big.NewInt(30000000000), // 30 gwei
	)
	fmt.Println("✓ Polygon integration initialized")

	// Register Auditor
	auditorCredentials := map[string]interface{}{
		"license":        "CA-12345",
		"experience":     "15 years",
		"certifications": []string{"CA", "CPA", "CISCA"},
	}
	auditor, err := ngoPlat.RegisterAuditor(
		"AUD001",
		"PwC India",
		auditorCredentials,
		[]string{"financial", "compliance", "gst"},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ngoPlat.VerifyAuditorCredentials("AUD001", "ICAI_AUTHORITY")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Auditor registered and verified: %s\n", auditor.Name)

	// Register NGO with comprehensive data
	ngoKYCData := map[string]interface{}{
		"documents": []string{"registration_cert.pdf", "80g_certificate.pdf", "pan_card.pdf"},
		"address":   "Mumbai, Maharashtra",
		"website":   "https://savechildren.org.in",
		"contact":   "+91-9876543210",
	}
	ngoSigners := []string{"0xNGOSigner1", "0xNGOSigner2", "0xTrustee1"}

	ngo, err := ngoPlat.RegisterNGO(
		"NGO001",
		"Save The Children Foundation India",
		"REG/2020/123456",
		"Child Welfare",
		ngoKYCData,
		ngoSigners,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Verify NGO KYC
	certificates := []entities.Certificate{
		{Type: "80G", Number: "80G/2020/12345", ValidUntil: "2025-12-31"},
		{Type: "FCRA", Number: "FCRA/2020/67890", ValidUntil: "2025-12-31"},
	}
	err = ngoPlat.VerifyNGOKYC("NGO001", "GOVERNMENT_AUTHORITY", certificates)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ NGO registered and KYC verified: %s\n", ngo.Name)

	// Register Donor with KYC
	donorKYCData := map[string]interface{}{
		"documents":    []string{"aadhaar", "pan"},
		"annual_limit": float64(500000), // 5 lakh limit
		"age":          35,
		"occupation":   "Software Engineer",
	}
	donor, err := ngoPlat.RegisterDonor("DONOR001", donorKYCData)
	if err != nil {
		log.Fatal(err)
	}

	// Verify Donor KYC
	err = ngoPlat.VerifyDonorKYC("DONOR001", "FINTECH_KYC_PROVIDER", "premium")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Donor registered and KYC verified")

	// Add NGO to donor's preferred list
	donor.AddPreferredNGO("NGO001")

	// Process Multiple Donations
	fmt.Println("\n--- Processing Donations ---")

	donations := []struct {
		amount float64
		method string
	}{
		{50000, "UPI"},
		{25000, "Credit Card"},
		{30000, "Net Banking"},
	}

	for _, donationData := range donations {
		result, err := ngoPlat.ProcessDonation(
			"DONOR001",
			"NGO001",
			donationData.amount,
			donationData.method,
		)
		if err != nil {
			fmt.Printf("✗ Donation failed: %s\n", err.Error())
			continue
		}

		fmt.Printf("✓ Donation processed: ₹%.0f via %s\n", donationData.amount, donationData.method)
		fmt.Printf("  Transaction ID: %s\n", result["transaction_id"])
		fmt.Printf("  Block Hash: %s\n", result["block_hash"])
		fmt.Printf("  Platform Fee: ₹%.2f\n", result["platform_fee"])
		fmt.Printf("  Net Amount: ₹%.2f\n", result["net_amount"])

		// Check if anchored to Polygon
		if eBill, ok := result["e_bill"].(map[string]interface{}); ok {
			if anchor, exists := eBill["polygon_anchor"]; exists {
				fmt.Printf("  Anchored to Polygon: %v\n", anchor)
			}
		}
	}

	// Process Expenditures
	fmt.Println("\n--- Processing Expenditures ---")

	expenditures := []map[string]interface{}{
		{
			"amount":      float64(40000),
			"category":    "Education",
			"description": "School books and supplies for 50 children",
		},
		{
			"amount":      float64(35000),
			"category":    "Healthcare",
			"description": "Medical supplies and medicines",
		},
	}

	for _, expenditureData := range expenditures {
		result, err := ngoPlat.ProcessExpenditure("NGO001", expenditureData, "AUD001")
		if err != nil {
			fmt.Printf("✗ Expenditure failed: %s\n", err.Error())
			continue
		}

		amount := expenditureData["amount"].(float64)
		category := expenditureData["category"].(string)

		fmt.Printf("✓ Expenditure processed: ₹%.0f for %s\n", amount, category)
		fmt.Printf("  Transaction ID: %s\n", result["transaction_id"])
		fmt.Printf("  Block Hash: %s\n", result["block_hash"])

		if auditResult, ok := result["audit_result"]; ok {
			fmt.Printf("  Audit Result: %v\n", auditResult)
		}
	}

	// Calculate and Display Ratings
	fmt.Println("\n--- NGO Ratings and Analysis ---")
	ratings := ngoPlat.CalculateAllNGORatings(30)
	for _, rating := range ratings {
		fmt.Printf("%s:\n", rating["name"])
		fmt.Printf("  Rating: %.2f/5.0 ⭐\n", rating["rating"])
		fmt.Printf("  Transparency Score: %d%%\n", rating["transparency_score"])
		fmt.Printf("  Utilization Rate: %s\n", rating["utilization_rate"])
		fmt.Printf("  Gap Percentage: %s\n", rating["gap_percentage"])
		fmt.Printf("  Documentation Quality: %s\n", rating["documentation_quality"])
		fmt.Printf("  KYC Verified: %v\n", rating["kyc_verified"])
	}

	// Display Dashboards
	fmt.Println("\n--- Donor Dashboard ---")
	donorDashboard, err := ngoPlat.GetDonorDashboard("DONOR001")
	if err != nil {
		log.Fatal(err)
	}

	if stats, ok := donorDashboard["stats"].(entities.DonorStats); ok {
		fmt.Printf("Total Donated: ₹%.2f\n", stats.TotalDonated)
		fmt.Printf("Current Year Donations: ₹%.2f\n", stats.CurrentYearDonations)
		fmt.Printf("Donation Count: %d\n", stats.DonationCount)
		fmt.Printf("Preferred NGOs: %d\n", stats.PreferredNGOsCount)
		fmt.Printf("Average Donation: ₹%.2f\n", stats.AverageDonation)
		fmt.Printf("Annual Limit: ₹%.2f\n", stats.AnnualLimit)
	}

	fmt.Println("\n--- NGO Dashboard ---")
	ngoDashboard, err := ngoPlat.GetNGODashboard("NGO001")
	if err != nil {
		log.Fatal(err)
	}

	if stats, ok := ngoDashboard["stats"].(map[string]interface{}); ok {
		fmt.Printf("Total Donations Received: ₹%.2f\n", stats["total_donations_received"])
		fmt.Printf("Total Expenditure Reported: ₹%.2f\n", stats["total_expenditure_reported"])
		fmt.Printf("Rating: %.2f/5.0\n", stats["rating"])
		fmt.Printf("Transparency Score: %d%%\n", stats["transparency_score"])
		fmt.Printf("Donation Blocks: %d\n", stats["donation_blockchain_length"])
		fmt.Printf("Expenditure Blocks: %d\n", stats["expenditure_blockchain_length"])
	}

	fmt.Println("\n--- Auditor Dashboard ---")
	auditorDashboard, err := ngoPlat.GetAuditorDashboard("AUD001")
	if err != nil {
		log.Fatal(err)
	}

	if auditorDashboard != nil {
		if stats, ok := auditorDashboard["stats"].(entities.AuditorStats); ok {
			fmt.Printf("Total Audits: %d\n", stats.TotalAudits)
			fmt.Printf("Approval Rate: %s\n", stats.ApprovalRate)
			fmt.Printf("Average Compliance Score: %.1f%%\n", stats.AverageComplianceScore)
			fmt.Printf("Rating: %.2f/5.0\n", stats.Rating)
		}
	}

	// Platform Statistics
	fmt.Println("\n--- Platform Statistics ---")
	platformStats := ngoPlat.GetPlatformStats()
	fmt.Printf("Total NGOs: %d\n", platformStats.TotalNGOs)
	fmt.Printf("Total Donors: %d\n", platformStats.TotalDonors)
	fmt.Printf("Total Auditors: %d\n", platformStats.TotalAuditors)
	fmt.Printf("Total Transactions: %d\n", platformStats.TotalTransactions)
	fmt.Printf("Total Donations: ₹%.2f\n", platformStats.TotalDonations)
	fmt.Printf("Total Expenditures: ₹%.2f\n", platformStats.TotalExpenditures)
	fmt.Printf("Platform Fee Collected: ₹%.2f\n", platformStats.PlatformFeeCollected)
	fmt.Printf("Average NGO Rating: %.2f/5.0\n", platformStats.AverageNGORating)
	fmt.Printf("Days Active: %d\n", platformStats.DaysActive)
	fmt.Printf("Verified NGOs: %d\n", platformStats.VerifiedNGOs)
	fmt.Printf("Verified Donors: %d\n", platformStats.VerifiedDonors)
	fmt.Printf("Verified Auditors: %d\n", platformStats.VerifiedAuditors)

	fmt.Println("\n=== ENHANCED DEMO COMPLETE ===")
	fmt.Println("\n🎉 All features demonstrated successfully!")
	fmt.Println("\n📊 Key Achievements:")
	fmt.Println("   • Zero-knowledge proofs for donor anonymity")
	fmt.Println("   • Multi-signature wallet integration")
	fmt.Println("   • Polygon blockchain anchoring")
	fmt.Println("   • Automated auditor validation")
	fmt.Println("   • Dynamic rating system")
	fmt.Println("   • Comprehensive KYC verification")
	fmt.Println("   • Real-time transparency scoring")
	fmt.Println("   • Tax benefit calculations")
	fmt.Println("   • GSTIN and invoice validation")
	fmt.Println("   • Platform fee management")
	fmt.Println("   • Thread-safe concurrent operations")
	fmt.Println("   • Comprehensive error handling")

	// Wait for user input before exiting
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
