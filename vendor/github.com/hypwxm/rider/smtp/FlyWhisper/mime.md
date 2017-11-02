From: from@qcode.co.uk
To: to@@qcode.co.uk
Subject: Example Email
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="MixedBoundaryString"

--MixedBoundaryString
Content-Type: multipart/related; boundary="RelatedBoundaryString"

--RelatedBoundaryString
Content-Type: multipart/alternative; boundary="AlternativeBoundaryString"

--AlternativeBoundaryString
Content-Type: text/plain;charset="utf-8"
Content-Transfer-Encoding: quoted-printable

This is the plain text part of the email.

--AlternativeBoundaryString
Content-Type: text/html;charset="utf-8"
Content-Transfer-Encoding: quoted-printable

<html>
  <body>=0D
    <img src=3D=22cid:masthead.png=40qcode.co.uk=22 width 800 height=3D80=
 =5C>=0D
    <p>This is the html part of the email.</p>=0D
    <img src=3D=22cid:logo.png=40qcode.co.uk=22 width 200 height=3D60 =5C=
>=0D
  </body>=0D
</html>=0D

--AlternativeBoundaryString--

--RelatedBoundaryString
Content-Type: image/jpgeg;name="logo.png"
Content-Transfer-Encoding: base64
Content-Disposition: inline;filename="logo.png"
Content-ID: <logo.png@qcode.co.uk>

amtsb2hiaXVvbHJueXZzNXQ2XHVmdGd5d2VoYmFmaGpremxidTh2b2hydHVqd255aHVpbnRyZnhu
dWkgb2l1b3NydGhpdXRvZ2hqdWlyb2h5dWd0aXJlaHN1aWhndXNpaHhidnVqZmtkeG5qaG5iZ3Vy
...
...
a25qbW9nNXRwbF0nemVycHpvemlnc3k5aDZqcm9wdHo7amlodDhpOTA4N3U5Nnkwb2tqMm9sd3An
LGZ2cDBbZWRzcm85eWo1Zmtsc2xrZ3g=

--RelatedBoundaryString
Content-Type: image/jpgeg;name="masthead.png"
Content-Transfer-Encoding: base64
Content-Disposition: inline;filename="masthead.png"
Content-ID: <masthead.png@qcode.co.uk>

aXR4ZGh5Yjd1OHk3MzQ4eXFndzhpYW9wO2tibHB6c2tqOTgwNXE0aW9qYWJ6aXBqOTBpcjl2MC1t
dGlmOTA0cW05dGkwbWk0OXQwYVttaXZvcnBhXGtsbGo7emt2c2pkZnI7Z2lwb2F1amdpNTh1NDlh
...
...
eXN6dWdoeXhiNzhuZzdnaHQ3eW9zemlqb2FqZWt0cmZ1eXZnamhka3JmdDg3aXV2dWd5aGVidXdz
dhyuhehe76YTGSFGA=

--RelatedBoundaryString--

--MixedBoundaryString
Content-Type: application/pdf;name="Invoice_1.pdf"
Content-Transfer-Encoding: base64
Content-Disposition: attachment;filename="Invoice_1.pdf"

aGZqZGtsZ3poZHVpeWZoemd2dXNoamRibngganZodWpyYWRuIHVqO0hmSjtyRVVPIEZSO05SVURF
SEx1aWhudWpoZ3h1XGh1c2loZWRma25kamlsXHpodXZpZmhkcnVsaGpnZmtsaGVqZ2xod2plZmdq
...
...
a2psajY1ZWxqanNveHV5ZXJ3NTQzYXRnZnJhZXdhcmV0eXRia2xhanNueXVpNjRvNWllc3l1c2lw
dWg4NTA0

--MixedBoundaryString
Content-Type: application/pdf;name="SpecialOffer.pdf"
Content-Transfer-Encoding: base64
Content-Disposition: attachment;filename="SpecialOffer.pdf"

aXBvY21odWl0dnI1dWk4OXdzNHU5NTgwcDN3YTt1OTQwc3U4NTk1dTg0dTV5OGlncHE1dW4zOTgw
cS0zNHU4NTk0eWI4OTcwdjg5MHE4cHV0O3BvYTt6dWI7dWlvenZ1em9pdW51dDlvdTg5YnE4N3Z3
...
...
OTViOHk5cDV3dTh5bnB3dWZ2OHQ5dTh2cHVpO2p2Ymd1eTg5MGg3ajY4bjZ2ODl1ZGlvcjQ1amts
dfnhgjdfihn=

--MixedBoundaryString--




Header
|From: email
|To: email
|MIME-Version: 1.0
|Content-Type: multipart/mixed; boundary="boundary1";
Message body
|multipart/mixed --boundary1
|--boundary1
|   multipart/related --boundary2
|   |--boundary2
|   |   multipart/alternative --boundary3
|   |   |--boundary3
|   |   |text/plain
|   |   |--boundary3
|   |   |text/html
|   |   |--boundary3--
|   |--boundary2
|   |Inline image
|   |--boundary2
|   |Inline image
|   |--boundary2--
|--boundary1
|Attachment1
|--boundary1
|Attachment2
|--boundary1
|Attachment3
|--boundary1--