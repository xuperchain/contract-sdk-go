var assert = require("assert");

var codePath = "../build/source_trace";
var lang = "go"
var type = "native"

function deploy(totalSupply) {
    return xchain.Deploy({
        name: "source_trace",
        code: codePath,
        lang: lang,
        type: type,
        init_args: {
            "admin": "XC1111111111111111@xuper"
        },
        options: { "account": "XC1111111111111111@xuper" }
    });
}

function CreateGoods() {
    c = deploy()
    resp = c.Invoke("CreateGoods", { "id": "id1", "desc": "goods1" })
    assert.equal(resp.Message, "missing initiator")
    resp = c.Invoke("CreateGoods", { "id": "id1", "desc": "goods1" }, { "account": "XC1111111111111111@xuper" })
    console.log(resp.Message)
    assert.equal(resp.Body, "id1")
    resp = c.Invoke("CreateGoods", { "id": "id2", "desc": "goods2" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "id2")
    {
        resp = c.Invoke("CreateGoods", { "id": "id2", "desc": "goods2" }, { "account": "XC1111111111111111@xuper" })
        assert.equal(resp.Message, "goods id2 already exists")
    }
    resp = c.Invoke("QueryRecords", { "id": "id1" })
    assert.equal(resp.Status, 200)
    assert.deepStrictEqual(JSON.parse(resp.Body), [{ "UpdateReccord": "0", "reason": "CREATE" }])
    resp = c.Invoke("UpdateGoods", { "id": "id1", "reason": "reason0" })
    assert.equal(resp.Message, "missing initiator")
    resp = c.Invoke("UpdateGoods", { "id": "id1", "reason": "reason0" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "1")
    {
        c.Invoke("UpdateGoods", { "id": "id1", "reason": "reason1" }, { "account": "XC1111111111111111@xuper" })
        c.Invoke("UpdateGoods", { "id": "id1", "reason": "reason2" }, { "account": "XC1111111111111111@xuper" })
    }
    resp = c.Invoke("QueryRecords", { "id": "id1" })
    // console.log(resp.Body)
    assert.deepStrictEqual(JSON.parse(resp.Body), [{ "UpdateReccord": "0", "reason": "CREATE" }, { "UpdateReccord": "1", "reason": "reason0" }, { "UpdateReccord": "2", "reason": "reason1" }, { "UpdateReccord": "3", "reason": "reason2" }])
}

function QueryRecords() {
    c = deploy()
    resp = c.Invoke("QueryRecords", { "id": "not_exist" })
    assert.deepStrictEqual("resp.Message", "goods not found")
}

Test("CreateGoods", CreateGoods)

