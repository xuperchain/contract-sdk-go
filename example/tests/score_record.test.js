var assert = require("assert");

var codePath = "../build/score_record";

var lang = "go"
var type = "native"

function deploy() {
    return xchain.Deploy({
        name: "award_manage",
        code: codePath,
        lang: lang,
        type: type,
        init_args: {
            "owner": "XC1111111111111111@xuper"
        },
        options: { "account": "XC1111111111111111@xuper" }
    });
}


function AddScore(t) {
    var c = deploy()
    var resp = c.Invoke("AddScore", { "user_id": "user1" })
    assert.equal(resp.Message, "missing initiator")
    var resp = c.Invoke("AddScore", { "user_id": "user1", "data": "data1" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "user1")
}

function QueryScore(t) {
    var c = deploy()
    resp = c.Invoke("AddScore", { "user_id": "user1", "data": "data1" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "user1")
    resp = c.Invoke("AddScore", { "user_id": "user2", "data": "data2" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "user2")

    resp = c.Invoke("AddScore", { "user_id": "user3" })
    assert.equal(resp.Message, "missing initiator")


    resp = c.Invoke("AddScore", { "user_id": "user3" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Status >= 500, true)
    console.log(resp.Message)
    // assert.equal(resp.Message, "missing data")


    resp = c.Invoke("QueryScore", { "user_id": "user1" })
    assert.equal(resp.Body, "data1")
}

function QueryOwner(t) {
    var c = deploy()
    var resp = c.Invoke("QueryOwner", {})
    assert.equal(resp.Body, "XC1111111111111111@xuper")
}


Test("QueryOwner", QueryOwner)
Test("QueryScore", QueryScore)
Test("AddScore", AddScore)