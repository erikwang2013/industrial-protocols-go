# परीक्षण रिपोर्ट — industrial-protocols-go

दिनांक: 2026-08-27
दायरा: Go मल्टी-मॉड्यूल वर्कस्पेस (go.work, 42 मॉड्यूल: kernel / protocols / examples)
कमांड: `go test ./...` और `go test -cover ./...` (क्योंकि वर्कस्पेस रूट में go.mod नहीं है, मॉड्यूल-दर-मॉड्यूल समतुल्य रूप से चलाया गया, परिणाम समान हैं)

## 1. समग्र निष्कर्ष

- **पूर्ण हरा**: सभी 42 मॉड्यूल और 52 टेस्ट-युक्त पैकेज पास (कई पूर्ण री-रन द्वारा सत्यापित, जिसमें समानांतर टेस्ट-लेखन अवधि भी शामिल है)।
- परीक्षण फ़ाइलें 106, परीक्षण फ़ंक्शन 659।
- `go vet ./...` — सभी मॉड्यूल में कोई चेतावनी नहीं।
- औसत स्टेटमेंट कवरेज **83.3%**; kernel कोर पैकेज कवरेज 95%~100%।

## 2. प्रति-मॉड्यूल परीक्षण आँकड़े और कवरेज

### kernel (11 पैकेज, examples को छोड़कर सर्वोच्च कवरेज)

| पैकेज | कवरेज |
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

### protocols के प्रमुख मॉड्यूल

| मॉड्यूल | कवरेज | मॉड्यूल | कवरेज |
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

> examples/modbus_basic में कोई परीक्षण फ़ाइल नहीं है (नमूना प्रोग्राम, अपेक्षित)।

## 3. सुधार सूची

### 3.1 सोर्स कोड बग फिक्स (tester-kernel / tester-protocols द्वारा खोजे और सत्यापित)

| फ़ाइल:पंक्ति | समस्या | समाधान |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` जिन नियमों में `Src`/`Dst`/`Map` nil है, उन पर `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` में nil डीरेफरेंस से पैनिक करता है | लूप की शुरुआत में अपूर्ण रूप से कॉन्फ़िगर किए गए नियमों को छोड़ें (मौजूदा "खराब नियम छोड़ें" शब्दार्थ के अनुरूप) |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` 1-बाइट अपवाद PDU (जैसे `0x81`) पर `pdu[1]` पढ़ते समय सीमा से बाहर पैनिक करता है | `len(pdu) < 2` जाँच जोड़ें, पार्सिंग त्रुटि लौटाएँ |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` जब बाइट काउंट (`pdu[1]`) वास्तविक डेटा से अधिक होता है तो स्लाइस सीमा से बाहर पैनिक (दुर्भावनापूर्ण/क्षतिग्रस्त फ्रेम) | `2+n > len(pdu)` सीमा जाँच जोड़ें, पार्सिंग त्रुटि लौटाएँ |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` मैसेज की लंबाई ऑफ़सेट 8 (प्रोटोकॉल संस्करण स्लॉट) पर वापस लिखता है, जबकि विनिर्देश ऑफ़सेट 4 माँगता है | `sizePos` (ऑफ़सेट 4) पहले से नोट करें, सही स्थान पर वापस लिखें |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel` का लंबाई फ़ील्ड कभी वापस नहीं लिखा जाता, हमेशा 0 रहता है | अंत में वास्तविक फ्रेम लंबाई के अनुसार वापस लिखें |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` फ्रेम को 16+len(cipReq) आवंटित करता है लेकिन CIP अनुरोध [28:] पर लिखा जाता है; अतिरिक्त 12 बाइट्स के कारण `copy` पूरा CIP अनुरोध चुपचाप छोड़ देता है; लंबाई फ़ील्ड भी 12 कम है | आवंटन 28+len(cipReq) करें, length फ़ील्ड 4+len(cipReq) करें |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` सीमा से बाहर `blockLen` पर स्लाइस सीमा से बाहर पैनिक करता है | `len(data) < 10+blockLen` जाँच जोड़ें |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode` की न्यूनतम लंबाई जाँच 4 है, लेकिन `unmarshalCAN` को `data[4:8]` चाहिए; 4~7 बाइट फ्रेम पैनिक करता है | न्यूनतम लंबाई 8 करें (वायर फॉर्मेट निश्चित 8 बाइट) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` `length<5` (`apduEnd < apduStart`) या अत्यधिक लंबाई (`apduEnd > len(data)`) पर स्लाइस सीमा से बाहर पैनिक करता है | `length < 5 || apduEnd > len(data)` की एक साथ जाँच करें और त्रुटि लौटाएँ |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP` का कार्य चर विनिर्देश के अनुसार `& 0xFF` से मास्क नहीं होता; ज्ञात वेक्टर `crc16DNP("123456789")` पर 0x69FF मिलता है (अपेक्षित 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode` की `len<5` जाँच fast-init सिंक फ्रेम की पहचान से पहले चलती है; 1-बाइट सिंक फ्रेम (0x55) गलत तरीके से अस्वीकार हो जाता है | पहले खाली फ्रेम और सिंक फ्रेम की जाँच करें, फिर लंबाई की |

### 3.2 परीक्षण सुधार (संकलन त्रुटियाँ / असर्शन त्रुटियाँ / डुप्लिकेट नाम / पुराने असर्शन ठीक किए गए)

| फ़ाइल:पंक्ति | समस्या | समाधान |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | स्थिरांक `0x68+0x10+0xF1+0x01`(362) को byte में असाइन करने पर संकलन विफल | अपेक्षा को चेकसम के निचले 8 बिट `0x6A` में बदलें |
| protocols/automotive/saej1850/saej1850_extra_test.go | `kernel.Codec` इंटरफ़ेस पर अनएक्सपोर्टेड मेथड/फ़ील्ड कॉल करने पर संकलन विफल | `newCodec` helper जोड़कर `*j1850Codec` होने की असर्शन करें; `TestEncodeMode01PID` का असर्शन ऑफ़सेट raw[4..6] (ID 4 बाइट + डेटा के पहले 4 बाइट) करें |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` मौजूदा परीक्षण से डुप्लिकेट नाम, संकलन विफल | नाम बदलकर `TestEncodeReadDefaultsFiber` करें (fiber वेरिएंट कवरेज बनाए रखें) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` मौजूदा परीक्षण से डुप्लिकेट नाम, संकलन विफल | नाम बदलकर `TestEncodeDirectArcCmd` करें |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | अपेक्षा `write 0x1000 A` थी, वास्तविक `0A` है (`%X` बाइट के अनुसार निश्चित दो अंक, मौजूदा परीक्षण की `DEADBEEF` परंपरा के अनुरूप) | अपेक्षा `0A` करें |
| protocols/bridge/sercos/sercos_extra_test.go:62 | वैसा ही, अपेक्षा `1` को `01` होना चाहिए | अपेक्षा `01` करें |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | पुराने bug व्यवहार (लंबाई ऑफ़सेट 8 पर) का असर्शन, सोर्स फिक्स के बाद अमान्य | प्रोटोकॉल संस्करण फ़ील्ड 0 होने का असर्शन करें (लंबाई TestExtraHELMessageSizeAtOffset4 द्वारा कवर होती है) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" परीक्षण का असर्शन उल्टा है (`err != nil` पर विफल, जबकि पार्सिंग त्रुटि ही अपेक्षित परिणाम है) | `err == nil` होने पर विफल करें |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | वैसा ही, असर्शन उल्टा | `err == nil` होने पर विफल करें |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | वैसा ही ×2 | `err == nil` होने पर विफल करें |

> नोट: `kernel/security/tls_test.go` के `stubTransport` इंटरफ़ेस कार्यान्वयन समस्या और `cclink/cclinkie` के CRC अपेक्षित मान, परीक्षण इंजीनियर द्वारा समानांतर कार्य में स्वयं सुधारे गए, इस बार संशोधित नहीं किए गए।

## 4. शेष जोखिम

1. **कम कवरेज मॉड्यूल**: profibus (35.7%), lonworks / asinterface / foundationfieldbus / iolink (37.9%) — परीक्षण केवल कुछ पथ कवर करते हैं; आगे Decode/Encode शाखाएँ और अपवाद पथ जोड़ने की सिफारिश की जाती है।
2. **modbus RTU CRC सत्यापित नहीं**: `decodeRTU` CRC सत्यापित नहीं करता (परीक्षण `TestExtraRTUDecodesCorruptCRC` इस GAP को स्पष्ट रूप से प्रलेखित करता है और वर्तमान स्थिति में पास होता है)। वास्तविक डिवाइस के साथ इंटरऑपरेबिलिटी के लिए इसे जोड़ने की सिफारिश है।
3. **canopen वायर फॉर्मेट सीमा**: `marshalCAN` केवल पहले 4 बाइट डेटा ले जाता है (`TestExtraSDOWritePayloadLostOnWire` द्वारा प्रलेखित); SDO मल्टी-बाइट पेलोड लिखना खो जाता है।
4. **हार्डवेयर-निर्भर स्किप**: cpci / pci / vme बिना हार्डवेयर वाले वातावरण में `t.Skip` करते हैं (उचित)।
5. **रिपॉज़िटरी में अस्थायी फ़ाइलें शेष**: रूट में अनट्रैक्ड `main.go` (प्रोब प्रोग्राम, गैर-मौजूद `probe/` मॉड्यूल का संदर्भ देता है), `scripts/`, `docs/*.png` — इस डिलीवरी का हिस्सा नहीं हैं, सफाई की सिफारिश की जाती है।
6. **परीक्षण अभी भी समानांतर में लिखे जा रहे हैं**: यह रिपोर्ट अंतिम पूर्ण हरे-परीक्षण स्नैपशॉट पर आधारित है; यदि परीक्षण इंजीनियर आगे `*_extra_test.go` फ़ाइलें जमा करते हैं, तो पूर्ण रिग्रेशन दोबारा चलाना आवश्यक होगा।
