import { keygen, build } from "./gobl.js";

window.gobl = {};
window.gobl.keygen = keygen;
window.gobl.build = build;

// const result = await keygen();
// console.log(`RESULT: ${result}`);

// try {
//     const buildResult = await build({
//         "data": {},
//         "privateykey": {},
//     });
//     console.log(`BUILD RESULT: ${build_result}`);
// } catch (e) {
//     console.log("BUILD ERROR: " + e)
// };

// const result2 = await keygen();
// console.log(`RESULT2: ${result2}`);

let goblData = {};

const exampleInputs = {};
exampleInputs.empty = `{
  "data": {},
  "privateykey": {}
}`;
exampleInputs.success = `{
	"$schema": "https://gobl.org/draft-0/envelope",
	"head": {
		"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
		"dig": {
			"alg": "sha256",
			"val": "579ae6960ff82e47a5478f8ed41728dba6298cab985cfea8caf695388fe0d721"
		}
	},
	"doc": {
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"region": "ES",
		"uuid": "3d7fdbdc-d037-11eb-a068-3e7e00ce5635",
		"code": "INV2021-001",
		"currency": "EUR",
		"issue_date": "2021-06-16",
		"supplier": {
			"tax_id": {
				"country": "ES",
				"code": "B91983379"
			},
			"name": "A Company Name S.L.",
			"people": [
				{
					"name": {
						"alias": "Paco",
						"given": "Francisco",
						"surname": "Smith"
					}
				}
			],
			"addresses": [
				{
					"num": "10",
					"street": "Calle Mayor",
					"locality": "Madrid",
					"region": "Madrid",
					"code": "28003",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "contact@company.com"
				}
			],
			"telephones": [
				{
					"label": "mobile",
					"num": "+34644123123"
				}
			]
		},
		"customer": {
			"tax_id": {
				"country": "ES",
				"code": "B-85905495"
			},
			"name": "Autofiscal S.L.",
			"addresses": [
				{
					"num": "1629",
					"street": "Calle Diseminado",
					"locality": "Miraflores de la Sierra",
					"region": "Madrid",
					"code": "28792",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "sam.lown@invopop.com"
				}
			]
		},
		"lines": [
			{
				"i": 1,
				"quantity": "20",
				"item": {
					"name": "Development services day rate",
					"price": "200.00"
				},
				"sum": "4000.00",
				"discount": {
					"value": "200.00",
					"reason": "just because"
				},
				"taxes": [
					{
						"cat": "VAT",
						"code": "STD"
					},
					{
						"cat": "IRPF",
						"code": "STD"
					}
				],
				"total": "3800.00"
			},
			{
				"i": 2,
				"quantity": "10",
				"item": {
					"name": "Something random",
					"price": "100.00"
				},
				"sum": "1000.00",
				"taxes": [
					{
						"cat": "VAT",
						"code": "STD"
					}
				],
				"total": "1000.00"
			},
			{
				"i": 3,
				"quantity": "5",
				"item": {
					"name": "Additional random services",
					"price": "34.50"
				},
				"sum": "172.50",
				"taxes": [
					{
						"cat": "VAT",
						"code": "RED"
					}
				],
				"total": "172.50"
			},
			{
				"i": 4,
				"quantity": "3",
				"item": {
					"name": "Impuesto local",
					"price": "5.00"
				},
				"sum": "15.00",
				"total": "15.00"
			}
		],
		"outlays": [
			{
				"i": 1,
				"desc": "Something paid for by us",
				"paid": "200.00"
			}
		],
		"totals": {
			"sum": "5187.50",
			"discount": "200.00",
			"total": "4987.50",
			"taxes": {
				"categories": [
					{
						"code": "VAT",
						"rates": [
							{
								"code": "STD",
								"base": "4800.00",
								"percent": "21.0%",
								"value": "1008.00"
							},
							{
								"code": "RED",
								"base": "172.50",
								"percent": "10.0%",
								"value": "17.25"
							}
						],
						"base": "4972.50",
						"value": "1025.25"
					},
					{
						"code": "IRPF",
						"retained": true,
						"rates": [
							{
								"code": "STD",
								"base": "3800.00",
								"percent": "15.0%",
								"value": "570.00"
							}
						],
						"base": "3800.00",
						"value": "570.00"
					}
				],
				"sum": "455.25"
			},
			"outlays": "200.00",
			"payable": "5642.75"
		},
		"payment": {
			"terms": {
				"code": "instant"
			},
			"instructions": {
				"code": "credit_transfer",
				"credit_transfer": [
					{
						"iban": "ES06 0128 0011 3901 0008 1391",
						"name": "Bankinter"
					}
				]
			}
		}
	}
}`

exampleInputs.testBuildSuccess = `{
	"$schema": "https://gobl.org/draft-0/envelope",
	"head": {
		"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
		"dig": {
			"alg": "sha256",
			"val": "00a1c7bb485818ec361c2ae243c78bdc9ce90ec638eced71c97df04325d84e09"
		}
	},
	"doc": {
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"uuid": "3d7fdbdc-d037-11eb-a068-3e7e00ce5635",
		"region": "ES",
		"code": "INV2021-001",
		"currency": "EUR",
		"issue_date": "2021-06-16",
		"supplier": {
			"tax_id": {
				"country": "ES",
				"code": "B91983379"
			},
			"name": "A Company Name S.L.",
			"people": [
				{
					"name": {
						"alias": "Paco",
						"given": "Francisco",
						"surname": "Smith"
					}
				}
			],
			"addresses": [
				{
					"num": "10",
					"street": "Calle Mayor",
					"locality": "Madrid",
					"region": "Madrid",
					"code": "28003",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "contact@company.com"
				}
			],
			"telephones": [
				{
					"label": "mobile",
					"num": "+34644123123"
				}
			]
		},
		"customer": {
			"tax_id": {
				"country": "ES",
				"code": "B-85905495"
			},
			"name": "Autofiscal S.L.",
			"addresses": [
				{
					"num": "1629",
					"street": "Calle Diseminado",
					"locality": "Miraflores de la Sierra",
					"region": "Madrid",
					"code": "28792",
					"country": "ES"
				}
			],
			"emails": [
				{
					"addr": "sam.lown@invopop.com"
				}
			]
		},
		"lines": [
			{
				"i": 1,
				"quantity": "20",
				"item": {
					"name": "Development services day rate",
					"price": "200.00"
				},
				"sum": "4000.00",
				"taxes": [
					{
						"cat": "VAT",
						"code": "STD"
					},
					{
						"cat": "IRPF",
						"code": "STD"
					}
				],
				"total": "4000.00"
			},
			{
				"i": 2,
				"quantity": "10",
				"item": {
					"name": "Something random",
					"price": "100.00"
				},
				"sum": "1000.00",
				"taxes": [
					{
						"cat": "VAT",
						"code": "STD"
					}
				],
				"total": "1000.00"
			},
			{
				"i": 3,
				"quantity": "5",
				"item": {
					"name": "Additional random services",
					"price": "34.50"
				},
				"sum": "172.50",
				"taxes": [
					{
						"cat": "VAT",
						"code": "RED"
					}
				],
				"total": "172.50"
			},
			{
				"i": 4,
				"quantity": "3",
				"item": {
					"name": "Impuesto local",
					"price": "5.00"
				},
				"sum": "15.00",
				"total": "15.00"
			}
		],
		"outlays": [
			{
				"i": 1,
				"desc": "Something paid for by us",
				"amount": "0"
			}
		],
		"totals": {
			"sum": "5187.50",
			"total": "5187.50",
			"taxes": {
				"categories": [
					{
						"code": "VAT",
						"rates": [
							{
								"code": "STD",
								"base": "5000.00",
								"percent": "21.0%",
								"amount": "1050.00"
							},
							{
								"code": "RED",
								"base": "172.50",
								"percent": "10.0%",
								"amount": "17.25"
							}
						],
						"base": "5172.50",
						"amount": "1067.25"
					},
					{
						"code": "IRPF",
						"retained": true,
						"rates": [
							{
								"code": "STD",
								"base": "4000.00",
								"percent": "15.0%",
								"amount": "600.00"
							}
						],
						"base": "4000.00",
						"amount": "600.00"
					}
				],
				"sum": "467.25"
			},
			"total_with_tax": "5654.75",
			"outlays": "0.00",
			"payable": "5654.75"
		},
		"payment": {
			"terms": {
				"code": "instant"
			},
			"instructions": {
				"code": "credit_transfer",
				"credit_transfer": [
					{
						"iban": "ES06 0128 0011 3901 0008 1391",
						"name": "Bankinter"
					}
				]
			}
		}
	},
	"sigs": ["sig data"]
}`

const testYaml = `title: "Test Message"
content: |-
  We hope you like this test message!`;

const testExampleFromWebsite = `{
  "$schema": "https://gobl.org/draft-0/envelope",
  "head": {
    "uuid": "8e69dd09-9adb-11ec-82ed-665181255c0a",
    "dig": {
      "alg": "sha256",
      "val": "45ac3115c8569a1789e58af8d0dc91ef3baa1fb71daaf38f5aef94f82b4d0033"
    },
    "draft": true
  },
  "doc": {
    "$schema": "https://gobl.org/draft-0/note/message",
    "title": "Test Message",
    "content": "We hope you like this test message!"
  },
  "sigs": [
    "eyJhbGciOiJFUzI1NiIsImtpZCI6ImE2YzM5MjBkLThlMzMtNDE1OS1iMzM3LTllNzQ2MTcxNmRmMSJ9.eyJ1dWlkIjoiOGU2OWRkMDktOWFkYi0xMWVjLTgyZWQtNjY1MTgxMjU1YzBhIiwiZGlnIjp7ImFsZyI6InNoYTI1NiIsInZhbCI6IjQ1YWMzMTE1Yzg1NjlhMTc4OWU1OGFmOGQwZGM5MWVmM2JhYTFmYjcxZGFhZjM4ZjVhZWY5NGY4MmI0ZDAwMzMifX0.VV9LRGEVPoO-tnOS-j6ItUEvYNcaQ1CbwCMN3qJorZXV3ON51wzalRuzJxulPnlFPtohWd_gc2Mf81MDIAK47Q"
  ]
}`;

const loadKey = async () => {
    const key = await keygen();
    goblData.key = key;
    document.getElementById("key").innerHTML = key;
}

const displayExample = async () => {
    document.getElementById("input-file").innerHTML = exampleInputs.success;
}

const processInput = async () => {
    const inputFile = document.getElementById("input-file").innerHTML;
    // console.log(inputFile)
    // const inputJSON = JSON.parse(inputFile);
    // const lol = {
    //     data: inputJSON,
    //     privatekey: goblData.key,
    //     sigs: []
    // }
    const buildResult = await build(JSON.parse(testExampleFromWebsite));

    document.getElementById("output-file").innerHTML = buildResult;
}

loadKey();
displayExample();
processInput();