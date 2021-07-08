var assert = require("assert");

var codePath = "../build/game_assets";
var lang = "go"
var type = "native"
function deploy(totalSupply) {
    return xchain.Deploy({
        name: "award_manage",
        code: codePath,
        lang: lang,
        type: type,
        init_args: {
            "admin": "XC1111111111111111@xuper"
        },
        options: { "account": "XC1111111111111111@xuper" }
    });
}


function AddAsset(t) {
    var c = deploy()
    resp = c.Invoke("AddAssetType", {
        "type_id": "type_id1",
        "type_desc": "type_desc1"
    }, { "account": "anonymous" })
    assert.equal(resp.Message, "you do not have permission to call this method")
    resp = c.Invoke("AddAssetType", { "type_id": "type_id1", "type_desc": "type_desc1" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "type_id1")
    resp = c.Invoke("AddAssetType", {
        "type_id": "type_id1",
        "type_desc": "type_desc1"
    }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Message, "asset type type_id1 already exists")
    return c
}
function ListAssetType(t) {
    var c = deploy(0)
    c.Invoke("AddAssetType", { "type_id": "type_id1", "type_desc": "type_desc1" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id2", "type_desc": "type_desc2" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id3", "type_desc": "type_desc3" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id4", "type_desc": "type_desc4" }, { "account": "XC1111111111111111@xuper" })
    resp = c.Invoke("ListAssetType", {})
    console.log(resp.Body)
}

function AssetOperations() {
    var c = deploy()
    resp = c.Invoke("AddAssetType", { "type_id": "type_id1", "type_desc": "type_desc1" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id2", "type_desc": "type_desc2" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id3", "type_desc": "type_desc3" }, { "account": "XC1111111111111111@xuper" })
    c.Invoke("AddAssetType", { "type_id": "type_id4", "type_desc": "type_desc4" }, { "account": "XC1111111111111111@xuper" })


    resp = c.Invoke("NewAssetToUser", { "user_id": "user_id1", "type_id": "type_id1", "asset_id": "asset_id1" })
    assert.equal(resp.Message, "missing initiator")
    {
        resp = c.Invoke("NewAssetToUser", { "user_id": "user_id1", "type_id": "type_id1", "asset_id": "asset_id1" }, { "account": "unknown" })
        assert.equal(resp.Message, "you do not have permission to call this method")

        resp = c.Invoke("NewAssetToUser", { "user_id": "user_id1", "type_id": "type_id100", "asset_id": "asset_id1" }, { "account": "XC1111111111111111@xuper" })
        assert.equal(resp.Message, "asset type type_id100 not found")
    }

    resp = c.Invoke("NewAssetToUser", { "user_id": "user_id1", "type_id": "type_id1", "asset_id": "asset_id1" }, { "account": "XC1111111111111111@xuper" })
    assert.equal(resp.Body, "asset_id1")

    {
        resp = c.Invoke("NewAssetToUser", { "user_id": "user_id2", "type_id": "type_id1", "asset_id": "asset_id1" }, { "account": "XC1111111111111111@xuper" })
        console.log(resp.Body)
        assert.equal(resp.Message, "asset asset_id1 exists")
    }
    {
        c.Invoke("NewAssetToUser", { "user_id": "user_id2", "type_id": "type_id3", "asset_id": "asset_id2" }, { "account": "XC1111111111111111@xuper" })

        c.Invoke("NewAssetToUser", { "user_id": "user_id2", "type_id": "type_id1", "asset_id": "asset_id3" }, { "account": "XC1111111111111111@xuper" })
        c.Invoke("NewAssetToUser", { "user_id": "user_id1", "type_id": "type_id1", "asset_id": "asset_id4" }, { "account": "XC1111111111111111@xuper" })
    }

    resp = c.Invoke("GetAssetsByUser", { "user_id": "user_id1" }, { "account": "XC1111111111111111@xuper" })
    console.log(resp.Body)
    assert.deepEqual(JSON.parse(resp.Body), [{ "id": "asset_id1", "type_id": "type_id1", "asset_desc": "type_desc1" }, { "id": "asset_id4", "type_id": "type_id1", "asset_desc": "type_desc1" }])



    resp = c.Invoke("GetAssetsByUser", { "user_id": "user_id2" }, { "account": "XC1111111111111111@xuper" })
    assert.deepEqual(JSON.parse(resp.Body), [{ "id": "asset_id2", "type_id": "type_id3", "asset_desc": "type_desc3" }, { "id": "asset_id3", "type_id": "type_id1", "asset_desc": "type_desc1" }])
    {
        resp = c.Invoke("TradeAsset", { "to": "user_id2", "asset_id": "asset_id2" }, { "account": "user_id1" })
        assert.equal(resp.Message, "asset asset_id2 of user user_id1 not found")
    }

    resp = c.Invoke("TradeAsset", { "to": "user_id2", "asset_id": "asset_id1" }, { "account": "user_id1" })
    assert.equal(resp.Status, 200)
    resp = c.Invoke("GetAssetsByUser", { "user_id": "user_id1" }, { "account": "XC1111111111111111@xuper" })
    assert.deepEqual(JSON.parse(resp.Body), [{ "id": "asset_id4", "type_id": "type_id1", "asset_desc": "type_desc1" }])

    resp = c.Invoke("GetAssetsByUser", { "user_id": "user_id2" }, { "account": "XC1111111111111111@xuper" })
    console.log(resp.Body)
    assert.deepEqual(JSON.parse(resp.Body), [{
        "id": "asset_id1", "type_id": "type_id1", "asset_desc": "type_desc1"
    }, {
        "id": "asset_id2", "type_id": "type_id3", "asset_desc": "type_desc3"
    }, {
        "id": "asset_id3", "type_id": "type_id1", "asset_desc": "type_desc1"
    }])
    {
        c.Invoke("NewAssetToUser", {
            "user_id": "XC1111111111111111@xuper", "type_id": "type_id1", "asset_id": "asset_id5"
        }, { "account": "XC1111111111111111@xuper" })
        resp = c.Invoke("GetAssetsByUser", {}, { "account": "XC1111111111111111@xuper" })
        assert.deepEqual(JSON.parse(resp.Body), [{
            "id": "asset_id5",
            "type_id": "type_id1",
            "asset_desc": "type_desc1"
        }])
    }
}


Test("AddAsset", AddAsset)
Test("ListAssetType", ListAssetType)
Test("AssetOperations", AssetOperations)
