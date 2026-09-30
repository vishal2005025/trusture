Backend
```
cd /Users/shanks/Documents/Projects/trusture
go run cmd/main.go
=== Enhanced NGO Transparency Platform Demo ===
✓ Platform initialized
✓ Polygon integration initialized
✓ Auditor registered and verified: PwC India
✓ NGO registered and KYC verified: Save The Children Foundation India
✓ Donor registered and KYC verified

--- Processing Donations ---
Block mined: 00a0fa284cdb816fb981dd46f0553127c9c651ce8d46829e4957bafb6ec60d0b
✓ Donation processed: ₹50000 via UPI
  Transaction ID: ddbe677a1b3edcb2a4ba9611132513e4
  Block Hash: 00a0fa284cdb816fb981dd46f0553127c9c651ce8d46829e4957bafb6ec60d0b
  Platform Fee: ₹500.00
  Net Amount: ₹49500.00
  Anchored to Polygon: &{0x2844cd65d9710410476ee84976ab875dae2704eb4a88e4380aa215e229cd7c26 3f5748a8ef4b55233b579be614b04cce8b5043cda7021c33f1120b5e8236b5ea 2025-11-10 04:46:11.698841 +0530 IST m=+0.202627876 50426887 23552 12}
Block mined: 001c456120b62770e488d85f6ce9bb8839c36f0ab9cdb8a2c2811f629956c15b
✓ Donation processed: ₹25000 via Credit Card
  Transaction ID: abb3285a0b28ac30ed8eeb1bba8a6924
  Block Hash: 001c456120b62770e488d85f6ce9bb8839c36f0ab9cdb8a2c2811f629956c15b
  Platform Fee: ₹250.00
  Net Amount: ₹24750.00
  Anchored to Polygon: &{0x1ee0acd226799718ea41ae42a6da7392958b84b58399aa032e61061fc718925c 521a4588af1c2238aa9c7a52709694ca80fa9ac18bf514bc967f0e0273245c82 2025-11-10 04:46:11.900355 +0530 IST m=+0.403995001 50690833 42350 12}
Block mined: 00f0bde73d2b4d5970ea723f46f1d2942adda7f804c8ebc267d5eeaff6bdd992
✓ Donation processed: ₹30000 via Net Banking
  Transaction ID: ad4907e7867ee6958f4c06a71764ec6a
  Block Hash: 00f0bde73d2b4d5970ea723f46f1d2942adda7f804c8ebc267d5eeaff6bdd992
  Platform Fee: ₹300.00
  Net Amount: ₹29700.00
  Anchored to Polygon: &{0xec5d7a4660e70a7382e6b8ffa0978725eeeaa73f3527515cf73d6acbad85fa9b 6a61d3062fc869c786ce0beb001d2e645bb05f33615a15f3d1a34ebbc9c03726 2025-11-10 04:46:12.102232 +0530 IST m=+0.605723793 50310904 46569 12}

--- Processing Expenditures ---
Block mined: 002752e351344708da4bca1df1ee6c5057c339a0cfdc872b46e040e167297747
✓ Expenditure processed: ₹40000 for Education
  Transaction ID: 355698c7e1db4d276c561dede35be98f
  Block Hash: 002752e351344708da4bca1df1ee6c5057c339a0cfdc872b46e040e167297747
  Audit Result: &{AUD-1762730172-AUD001 355698c7e1db4d276c561dede35be98f AUD001 2025-11-10 04:46:12.102332 +0530 IST m=+0.605824001 85 [Missing payment proof] Approve - Good compliance with minor observations  5f28ad482c6c57f85aac0d2d304dd7016af759c739ac885c840df1365e5ea87c}
Block mined: 00cbc91ceea7c9b395b8732a99b3cd4feb337603f51fab33f9b7ab2cab0668f1
✓ Expenditure processed: ₹35000 for Healthcare
  Transaction ID: fe08e062f58f867ae267e2fa4925e04a
  Block Hash: 00cbc91ceea7c9b395b8732a99b3cd4feb337603f51fab33f9b7ab2cab0668f1
  Audit Result: &{AUD-1762730172-AUD001 fe08e062f58f867ae267e2fa4925e04a AUD001 2025-11-10 04:46:12.305201 +0530 IST m=+0.808551793 85 [Missing payment proof] Approve - Good compliance with minor observations  622149df507275fb8ee57ab0e866ed98e7cabbb0808c281543f4c5156ef9613a}

--- NGO Ratings and Analysis ---
Save The Children Foundation India:
  Rating: 5.00/5.0 ⭐
  Transparency Score: 100%
  Utilization Rate: 72.15%
  Gap Percentage: 27.85%
  Documentation Quality: 30.0%
  KYC Verified: true

--- Donor Dashboard ---
Total Donated: ₹103950.00
Current Year Donations: ₹103950.00
Donation Count: 3
Preferred NGOs: 1
Average Donation: ₹34650.00
Annual Limit: ₹5000000.00

--- NGO Dashboard ---
Total Donations Received: ₹103950.00
Total Expenditure Reported: ₹75000.00
Rating: 5.00/5.0
Transparency Score: 100%
Donation Blocks: 4
Expenditure Blocks: 3

--- Auditor Dashboard ---
Total Audits: 2
Approval Rate: 100.00%
Average Compliance Score: 85.0%
Rating: 5.00/5.0

--- Platform Statistics ---
Total NGOs: 1
Total Donors: 1
Total Auditors: 1
Total Transactions: 5
Total Donations: ₹103950.00
Total Expenditures: ₹75000.00
Platform Fee Collected: ₹1039.50
Average NGO Rating: 5.00/5.0
Days Active: 0
Verified NGOs: 1
Verified Donors: 1
Verified Auditors: 1

=== ENHANCED DEMO COMPLETE ===

🎉 All features demonstrated successfully!

📊 Key Achievements:
   • Zero-knowledge proofs for donor anonymity
   • Multi-signature wallet integration
   • Polygon blockchain anchoring
   • Automated auditor validation
   • Dynamic rating system
   • Comprehensive KYC verification
   • Real-time transparency scoring
   • Tax benefit calculations
   • GSTIN and invoice validation
   • Platform fee management
   • Thread-safe concurrent operations
   • Comprehensive error handling

Press Enter to exit...
^Csignal: interrupt
➜  trusture git:(main) ✗ 
```




Front end


```
cd /Users/shanks/Documents/Projects/trusture/Frontend
npm install
npm run dev

added 243 packages, and audited 244 packages in 3s

58 packages are looking for funding
  run `npm fund` for details

2 moderate severity vulnerabilities

To address all issues, run:
  npm audit fix

Run `npm audit` for details.

> frontend@0.0.0 dev
> vite


  VITE v5.4.10  ready in 908 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
 ^C
(base) ➜  Frontend git:(main) ✗ 



```





