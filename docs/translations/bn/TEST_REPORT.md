# টেস্ট রিপোর্ট — industrial-protocols-go

তারিখ: 2026-08-27
পরিধি: Go মাল্টি-মডিউল ওয়ার্কস্পেস (go.work, ৪২টি মডিউল: kernel / protocols / examples)
কমান্ড: `go test ./...` এবং `go test -cover ./...` (ওয়ার্কস্পেস রুটে go.mod না থাকায় মডিউলভিত্তিক সমতুল্য নির্বাহ করা হয়েছে, ফলাফল একই)

## ১. সামগ্রিক উপসংহার

- **সম্পূর্ণ সবুজ**: ৪২টি মডিউলের ৫২টি টেস্টযুক্ত প্যাকেজ সব পাস করেছে (একাধিকবার সম্পূর্ণ রিরান দিয়ে যাচাই করা হয়েছে, সমান্তরাল টেস্ট লেখার সময়সহ)।
- টেস্ট ফাইল ১০৬টি, টেস্ট ফাংশন ৬৫৯টি।
- `go vet ./...` সব মডিউলে কোনো ওয়ার্নিং নেই।
- গড় স্টেটমেন্ট কভারেজ **৮৩.৩%**; kernel কোর প্যাকেজগুলোতে কভারেজ ৯৫%~১০০%।

## ২. মডিউলভিত্তিক টেস্ট পরিসংখ্যান ও কভারেজ

### kernel (১১টি প্যাকেজ, examples ছাড়া সর্বোচ্চ কভারেজ)

| প্যাকেজ | কভারেজ |
|---|---|
| kernel | 100% |
| kernel/bridge | 80.0% |
| kernel/config | 100% |
| kernel/connection | 95.3% |
| kernel/event | 100% |
| kernel/gateway | 100% |
| kernel/metrics | 100% |
| kernel/pipeline | 100% |
| kernel/security | 100% |
| kernel/transport | 100% |
| kernel/vendor | 100% |

### protocols গুরুত্বপূর্ণ মডিউল

| মডিউল | কভারেজ | মডিউল | কভারেজ |
|---|---|---|---|
| automotive/flexray | 79.4% | fieldbus/canopen | 73.9% |
| automotive/kline | 93.3% | fieldbus/cclink | 93.8% |
| automotive/lin | 92.5% | fieldbus/cclinkie | 73.5% |
| automotive/most | 100% | fieldbus/devicenet | 73.2% |
| automotive/saej1850 | 84.2% | fieldbus/dnp3 | 88.9% |
| bridge/controlnet | 100% | fieldbus/hart | 96.6% |
| bridge/ethercat | 100% | fieldbus/iec61850 | 91.8% |
| bridge/interbus | 70.0% | iot/hartip | 88.6% |
| bridge/isa100 | 95.2% | iot/mqtt | 89.1% |
| bridge/powerlink | 97.8% | ethernet/bacnet | 95.2% |
| bridge/sercos | 97.9% | ethernet/ethernetip | 96.9% |
| bridge/sercos1 | 95.6% | ethernet/modbus | 67.7% |
| bridge/safej1850 | 100% | ethernet/opcua | 95.8% |
| bridge/wirelesshart | 95.2% | ethernet/profinet | 95.4% |
| building/dali | 92.9% | system/cpci | 78.6% |
| building/lonworks | **37.9%** | system/pci | 78.6% |
| fieldbus/asinterface | **37.9%** | system/vme | 78.6% |
| fieldbus/foundationfieldbus | **37.9%** | fieldbus/profibus | **35.7%** |
| fieldbus/iolink | **37.9%** | bridge/lightbus / modbusplus / worldfip | 72~73.5% |

> examples/modbus_basic এ কোনো টেস্ট ফাইল নেই (একটি উদাহরণ প্রোগ্রাম, প্রত্যাশিত)।

## ৩. মেরামতের তালিকা

### ৩.১ সোর্স কোড বাগ মেরামত (tester-kernel / tester-protocols দ্বারা আবিষ্কৃত ও যাচাইকৃত)

| ফাইল:লাইন | সমস্যা | মেরামত |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform`-এ `Src`/`Dst`/`Map` nil থাকা রুলের ক্ষেত্রে `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` এ nil ডিরেফারেন্স panic | লুপের শুরুতে অসম্পূর্ণ কনফিগকৃত রুল স্কিপ করা (বিদ্যমান "খারাপ রুল স্কিপ" সেমান্টিক্সের সাথে সামঞ্জস্যপূর্ণ) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` ১-বাইট অস্বাভাবিক PDU-তে (যেমন `0x81`) `pdu[1]` পড়তে গিয়ে আউট-অফ-বাউন্ডস panic | `len(pdu) < 2` চেক যোগ করা, পার্সিং এরর রিটার্ন করে |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU`-এর বাইট কাউন্ট (`pdu[1]`) প্রকৃত ডেটার চেয়ে বড় হলে স্লাইস আউট-অফ-বাউন্ডস panic (দূষিত/ক্ষতিকারক ফ্রেম) | `2+n > len(pdu)` বাউন্ডারি চেক যোগ করা, পার্সিং এরর রিটার্ন করে |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` মেসেজ দৈর্ঘ্য অফসেট 8 (প্রোটোকল ভার্সন স্লট)-এ ব্যাকফিল করে; স্পেক অনুযায়ী অফসেট 4 হওয়া উচিত | `sizePos` (অফসেট 4) আগে থেকে রেকর্ড করে সঠিক অবস্থানে ব্যাকফিল করা |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel`-এর লেন্থ ফিল্ড কখনো ব্যাকফিল হয় না, সর্বদা 0 থাকে | শেষে প্রকৃত ফ্রেম দৈর্ঘ্য অনুযায়ী ব্যাকফিল করা |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` ফ্রেম বরাদ্দ করে 16+len(cipReq) কিন্তু CIP রিকোয়েস্ট লেখা হয় [28:]-এ; অতিরিক্ত ১২ বাইটের কারণে `copy` নীরবে পুরো CIP রিকোয়েস্ট বাদ দেয়; লেন্থ ফিল্ডও ১২ কম | বরাদ্দ 28+len(cipReq) করা, length ফিল্ড 4+len(cipReq) করা |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` সীমা-অতিরিক্ত `blockLen`-এ স্লাইস আউট-অফ-বাউন্ডস panic | `len(data) < 10+blockLen` চেক যোগ করা |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode`-এর ন্যূনতম দৈর্ঘ্য চেক 4, কিন্তু `unmarshalCAN`-এর জন্য `data[4:8]` প্রয়োজন; 4~7 বাইটের ফ্রেমে panic | ন্যূনতম দৈর্ঘ্য 8 করা (ওয়্যার ফরম্যাট নির্দিষ্ট ৮ বাইট) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU`-এ `length<5` (`apduEnd < apduStart`) বা অতিরিক্ত বড় `length` (`apduEnd > len(data)`) হলে স্লাইস আউট-অফ-বাউন্ডস panic | `length < 5 || apduEnd > len(data)` একীভূতভাবে যাচাই করে এরর রিটার্ন করা |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP`-এর ওয়ার্কিং ভেরিয়েবল স্পেক অনুযায়ী `& 0xFF` মাস্ক করা হয়নি; পরিচিত ভেক্টর `crc16DNP("123456789")` দেয় 0x69FF (হওয়া উচিত 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode`-এর `len<5` চেক fast-init সিঙ্ক ফ্রেম নির্ণয়ের আগে চলে; ১-বাইট সিঙ্ক ফ্রেম (0x55) ভুলভাবে প্রত্যাখ্যাত হয় | আগে খালি ফ্রেম ও সিঙ্ক ফ্রেম নির্ণয়, তারপর দৈর্ঘ্য যাচাই |

### ৩.২ টেস্ট মেরামত (কম্পাইল এরর / অ্যাসার্শন এরর / নামের দ্বন্দ্ব / অপ্রচলিত অ্যাসার্শন সংশোধন)

| ফাইল:লাইন | সমস্যা | মেরামত |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | কনস্ট্যান্ট `0x68+0x10+0xF1+0x01`(362) byte-তে অ্যাসাইন করলে কম্পাইল ব্যর্থ | প্রত্যাশা চেকসামের নিম্ন ৮ বিট `0x6A`-তে পরিবর্তিত |
| protocols/automotive/saej1850/saej1850_extra_test.go | `kernel.Codec` ইন্টারফেসে আনএক্সপোর্টেড মেথড/ফিল্ড কল করায় কম্পাইল ব্যর্থ | নতুন `newCodec` হেল্পার যোগ করে `*j1850Codec` হিসেবে অ্যাসার্শন; `TestEncodeMode01PID`-এর অফসেট সংশোধন করে raw[4..6] (ID ৪ বাইট + ডেটার প্রথম ৪ বাইট) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` বিদ্যমান টেস্টের সাথে নামের দ্বন্দ্বে কম্পাইল ব্যর্থ | নতুন নাম `TestEncodeReadDefaultsFiber` (fiber ভেরিয়েন্ট কভারেজ রেখে) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` বিদ্যমান টেস্টের সাথে নামের দ্বন্দ্বে কম্পাইল ব্যর্থ | নতুন নাম `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | প্রত্যাশা `write 0x1000 A`; প্রকৃত আউটপুট `0A` (`%X` বাইট অনুযায়ী নির্দিষ্ট দুই অঙ্ক, বিদ্যমান টেস্ট `DEADBEEF` কনভেনশনের সাথে সামঞ্জস্যপূর্ণ) | প্রত্যাশা `0A`-তে পরিবর্তিত |
| protocols/bridge/sercos/sercos_extra_test.go:62 | উপরের মতোই, প্রত্যাশা `1` হওয়া উচিত `01` | প্রত্যাশা `01`-তে পরিবর্তিত |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | পুরনো বাগের আচরণ অ্যাসার্শন করা হয়েছিল (দৈর্ঘ্য অফসেট 8-এ); সোর্স মেরামতের পর অবৈধ | প্রোটোকল ভার্সন ফিল্ড 0 কিনা অ্যাসার্শন করা (দৈর্ঘ্য TestExtraHELMessageSizeAtOffset4 কভার করে) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" টেস্টের অ্যাসার্শন উল্টানো (`err != nil` হলে ব্যর্থ, কিন্তু পার্সিং এররই প্রত্যাশিত ফলাফল) | `err == nil` হলে ব্যর্থ — এভাবে পরিবর্তিত |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | উপরের মতোই, অ্যাসার্শন উল্টানো | `err == nil` হলে ব্যর্থ — এভাবে পরিবর্তিত |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | উপরের মতোই ×২ | `err == nil` হলে ব্যর্থ — এভাবে পরিবর্তিত |

> নোট: `kernel/security/tls_test.go`-এর `stubTransport` ইন্টারফেস বাস্তবায়ন সমস্যা এবং `cclink/cclinkie`-এর CRC প্রত্যাশা মান টেস্ট ইঞ্জিনিয়ার সমান্তরাল কাজে নিজে সংশোধন করেছেন; এবার স্পর্শ করা হয়নি।

## ৪. অবশিষ্ট ঝুঁকি

1. **কম কভারেজ মডিউল**: profibus (35.7%), lonworks / asinterface / foundationfieldbus / iolink (37.9%) — টেস্ট শুধুমাত্র অল্প কয়েকটি পাথ কভার করে; পরে Decode/Encode ব্রাঞ্চ ও ব্যতিক্রম পাথ যোগ করার পরামর্শ।
2. **modbus RTU CRC যাচাই হয় না**: `decodeRTU` CRC যাচাই করে না (টেস্ট `TestExtraRTUDecodesCorruptCRC` এই GAP স্পষ্টভাবে ডকুমেন্ট করে এবং বর্তমান অবস্থায়ই পাস হয়)। প্রকৃত ডিভাইসের সাথে ইন্টারঅপারেবিলিটির জন্য যোগ করার পরামর্শ।
3. **canopen ওয়্যার ফরম্যাট সীমাবদ্ধতা**: `marshalCAN` শুধুমাত্র প্রথম ৪ বাইট ডেটা বহন করে (`TestExtraSDOWritePayloadLostOnWire` ডকুমেন্টেড); SDO দিয়ে মাল্টি-বাইট পেলোড লিখলে ডেটা হারিয়ে যাবে।
4. **হার্ডওয়্যার নির্ভরতা স্কিপ**: cpci / pci / vme হার্ডওয়্যারবিহীন পরিবেশে `t.Skip` করে (যৌক্তিক)।
5. **রিপোজিটরিতে অবশিষ্ট টেম্পোরারি ফাইল**: রুটে আনট্র্যাকড `main.go` (প্রোব প্রোগ্রাম, অবিদ্যমান `probe/` মডিউল রেফারেন্স করে), `scripts/`, `docs/*.png` — এই ডেলিভারির অংশ নয়, পরিষ্কার করার পরামর্শ।
6. **টেস্ট এখনও সমান্তরালে যুক্ত হচ্ছে**: এই রিপোর্ট শেষ সম্পূর্ণ সবুজ রানের স্ন্যাপশটের উপর ভিত্তি করে; টেস্ট ইঞ্জিনিয়ার পরে আরও `*_extra_test.go` জমা দিলে সম্পূর্ণ রিগ্রেশন পুনরায় চালাতে হবে।
