var assert = require("assert");

var codePath = "../build/features";

Test("caller", function (t) {
    var features1 = xchain.Deploy({
        name: "features1",
        code: codePath,
        lang: "go",
        init_args: {},
        type: "native",
        options: { "account": "XC1111111111111111@xuper" }
    })
    var features2 = xchain.Deploy({
        name: "features2",
        code: codePath,
        lang: "go",
        init_args: {},
        type: "native",
        options: { "account": "XC1111111111111111@xuper" }
    })
    resp = features1.Invoke("Call", {
        "contract": "features2",
        "method": "Caller"
    })
    console.log(resp.Message)
    assert.equal(resp.Body, "features1")
})