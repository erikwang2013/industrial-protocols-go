# PCI / PCIe प्रोटोकॉल SDK

Linux sysfs के माध्यम से PCI और PCI Express बस एक्सेस के लिए `kernel.Protocol` लागू करता है।

## प्रोटोकॉल

- **नाम**: `pci`
- **वेरिएंट**: `pci`
- **डिफ़ॉल्ट पोर्ट**: 0 (मेमोरी-मैप्ड बस)
- **ट्रांसपोर्ट**: sysfs `/sys/bus/pci/devices/<BDF>/config`

कोडेक एप्लिकेशन और PCI कॉन्फ़िगरेशन स्पेस के बीच कच्चे बाइट्स को पारित
करता है। रीड और राइट सीधे sysfs कॉन्फ़िगरेशन फ़ाइल के माध्यम से डिवाइस के
कॉन्फ़िगरेशन रजिस्टरों तक जाते हैं।

## कर्नेल आवश्यकताएँ

निम्नलिखित कर्नेल मॉड्यूल लोड होने चाहिए:

- `pcieport` -- PCI Express पोर्ट ड्राइवर
- `pci_sysfs` -- sysfs PCI इंटरफ़ेस (अधिकांश कर्नेल में निर्मित)

sysfs फाइलसिस्टम `/sys` पर माउंट होना चाहिए। यह सभी आधुनिक Linux
वितरणों पर डिफ़ॉल्ट है।

आवश्यक अनुमतियाँ:
- कॉन्फ़िगरेशन स्पेस एक्सेस के लिए रूट एक्सेस या `CAP_SYS_ADMIN`
- कॉन्फ़िगरेशन फ़ाइल root:root के स्वामित्व में mode 0600 के साथ

## ड्राइवर

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## BDF फॉर्मेट

बस पता एक BDF (Bus:Device.Function) स्ट्रिंग है:
- `0000:00:1f.3` -- domain 0000, bus 00, device 1f, function 3
- `0000:01:00.0` -- domain 0000, bus 01, device 00, function 0
