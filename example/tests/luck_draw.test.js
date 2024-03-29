var assert = require("assert");

var codePath = "../build/luck_draw";
var lang = "go"
var type = "native"

function deploy(totalSupply) {
    return xchain.Deploy({
        name: "luck_draw",
        code: codePath,
        lang: lang,
        type: type,
        init_args: {
            "admin": "XC1111111111111111@xuper"
        },
        options: { "account": "XC1111111111111111@xuper" },
    });
}


Test("LuckDraw", function (t) {
    c = deploy()
    c.Invoke("GetLuckId", {}, { "account": "user1" })
    c.Invoke("GetLuckId", {}, { "account": "user2" })
    c.Invoke("GetLuckId", {}, { "account": "user3" })
    c.Invoke("GetLuckId", {}, { "account": "user4" })
    resp = c.Invoke("GetLuckId", {}, { "account": "user5" })
    assert.equal(resp.Body, "5")

    resp = c.Invoke("GetLuckId", {}, { "account": "user1" })
    assert.equal(resp.Body, "1")

    resp = c.Invoke("StartLuckDraw", {}, { "account": "nobody" })
    assert.equal(resp.Message, "you do not have permission to call this method")

    resp = c.Invoke("StartLuckDraw", { "seed": "100" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Message, "")
    assert.equal(resp.Status, 200)
    resp = c.Invoke("GetResult", {})
    assert.equal(resp.Message, "")
    assert.equal(resp.Status, 200)

    resp = c.Invoke("GetLuckId", {}, { "account": "user5" })
    assert.equal(resp.Message, "the luck draw has finished")
})