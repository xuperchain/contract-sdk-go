var assert = require("assert");

var codePath = "../build/counter";

var lang = "go"
var type = "native"
function deploy() {
    return xchain.Deploy({
        name: "counter",
        code: codePath,
        lang: lang,
        type: type,
        init_args: {
            "creator": "XC1111111111111111@xuper"
        },
        options: { "account": "XC1111111111111111@xuper" }
    });
}

Test("Increase", function (t) {
    var c = deploy();
    var resp = c.Invoke("Increase", {
        "key": "XC1111111111111111@xuper"
    }, { "name": "11111" });
    assert.equal(resp.Body, "1");
})

Test("Get", function (t) {
    var c = deploy()
    c.Invoke("Increase", {
        "key": "XC1111111111111111@xuper"
    });
    var resp = c.Invoke("Get", { "key": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "1")
})