# MOST সিরিয়াল প্রোটোকল SDK

MOST (Media Oriented Systems Transport) একটি উচ্চ-গতির মাল্টিমিডিয়া নেটওয়ার্ক প্রযুক্তি, যা মূলত অটোমোটিভ ইনফোটেইনমেন্ট সিস্টেমে ব্যবহৃত হয়। এই প্যাকেজটি AT-কমান্ড ইন্টারফেসসহ একটি সিরিয়াল অ্যাডাপ্টারের উপর MOST কোডেক সরবরাহ করে।

## প্রোটোকল ওভারভিউ

MOST ফাইবার-অপটিক ফিজিক্যাল লেয়ারের উপর সিঙ্ক্রোনাস সিরিয়াল কমিউনিকেশন ব্যবহার করে। এই বাস্তবায়নটি 115200 বডে AT-কমান্ড ইন্টারফেস প্রদানকারী একটি সিরিয়াল অ্যাডাপ্টারের মাধ্যমে সংযুক্ত হয়।

### AT কমান্ড

| কমান্ড            | বিবরণ                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | ঠিকানা থেকে বাইট পড়া |
| `AT+WRITE=<addr>,<hex>` | ঠিকানায় হেক্স বাইট লেখা |
| `AT+STATUS`             | রিং/নেটওয়ার্ক স্ট্যাটাস জিজ্ঞাসা |

### রেসপন্স ফরম্যাট

- `+OK:<hex_data>` -- সফল রেসপন্স
- `+ERR:<code>` -- এরর রেসপন্স

## হার্ডওয়্যার প্রয়োজনীয়তা

- সঠিক টার্মিনেশনসহ MOST ফাইবার-অপটিক নেটওয়ার্ক
- MOST-থেকে-সিরিয়াল অ্যাডাপ্টার (যেমন MOST150 USB অ্যাডাপ্টার)
- 115200 বড, 8N1-এ সিরিয়াল পোর্ট

## ব্যবহার

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## ড্রাইভার

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## সমর্থিত ফাংশন

| ফাংশন | বিবরণ                 |
|----------|-----------------------------|
| `read`   | একটি MOST ঠিকানা থেকে পড়া    |
| `write`  | একটি MOST ঠিকানায় ডেটা লেখা |
| `status` | রিং/নেটওয়ার্ক স্ট্যাটাস জিজ্ঞাসা   |

## টেস্টিং

```bash
go test ./... -v
```
