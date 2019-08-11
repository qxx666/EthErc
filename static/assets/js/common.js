var headerpath = window.location.pathname;

$('.withDrawSearch').keydown(function (e) {
    if (e.keyCode == 13) {
        window.location.href = headerpath + "?address=" + $(".withDrawSearch").val();
    }
});

if (headerpath == "/admin/index.html") {
    var t1 = window.setInterval(summaryingTasks, 5000);
}

function showTx(tx) {
    layer.alert(tx)
}

function reWithdraw(id) {

    layer.prompt({title: '请输入主账户密码', formType: 1}, function(pass, index){

        $.ajax({
            type: "POST",
            url: '/admin/withdraw/reWithdraw.html',
            data: {
                withdrawId: id,
                password:pass
            },
            success: function (data) {
                alert(data.message);
                window.location.reload();
            }
        });

        layer.close(index);
    });


}

function summaryingTasks() {
    $.ajax({
        type: "GET",
        url: '/admin/summary/summarying.html',
        data: {},
        success: function (data) {
            if (data.code >= 0) {
                var html = "";
                jQuery.each(data.data, function (i, item) {
                    var percent = (item.finish_task_number / item.task_number) * 100
                    var liHtml = "<li class=\"am-list-border\">\n" +
                        "                                            <div class=\"cosA\">\n" +
                        "                                                <span class=\"am-table-bordered\">\n" +
                        "                                                    <i class=\"am-icon-tag\"></i>\n" +
                        "                                                    " + item.task_type + "汇总定时任务\n" +
                        "                                                </span>\n" +
                        "                                                &nbsp;\n" +
                        "                                                <span class=\"am-table-bordered\">\n" +
                        "                                                    <i class=\"am-icon-tag\"></i>\n" +
                        "                                                    开始时间:" + item.created_at + "\n" +
                        "                                                </span>\n" +
                        "                                                <div class=\"am-progress am-progress-striped\">\n" +
                        "                                                    <div class=\"am-progress-bar\" style=\"width: " + percent + "%\">" + percent + "%</div>\n" +
                        "                                                </div>\n" +
                        "                                                <span>\n" +
                        // "                                                    <a>暂停任务</a>\n" +
                        // "                                                     <a>结束任务</a>\n" +
                        "                                                     <a>查看详情</a>\n" +
                        "                                                </span>\n" +
                        "                                            </div>\n" +
                        "\n" +
                        "                                        </li>";
                    html = html + liHtml;
                });

                $(".task-ing-list").html(html);
            }
        }
    });
}

$("#coinSelector").change(function () {
    var coin = $(this).val();
    $.ajax({
        type: "POST",
        url: '/admin/account/setDefaultCoin.html',
        data: {
            coin: coin
        },
        success: function (data) {
            if (data.code >= 0) {
                window.location.reload();
            }
        }
    });

});

function updateSummarySetting() {
    $.ajax({
        type: "POST",
        url: '/admin/setting/summarySetting.html',
        data: $("#summarySettingForm").serializeArray(),
        success: function (data) {
            alert(data.message);
        }
    });
}

function setCurrentBlockNumber() {
    $.ajax({
        type: "POST",
        url: '/admin/setting/setCurrentBlockNumber.html',
        data: $("#currentBlockNumberForm").serializeArray(),
        success: function (data) {
            alert(data.message);
        }
    });
}

function rechargeStatus() {

    $.ajax({
        type: "POST",
        url: '/admin/changeRechargeStatus.html',
        data: {},
        success: function (data) {
            alert(data.message);
        }
    });
}

function summaryStatus() {

    $.ajax({
        type: "POST",
        url: '/admin/changeSummaryStatus.html',
        data: {},
        success: function (data) {
            alert(data.message);
        }
    });
}

function transferForm() {
    var index = layer.load(0, {shade: false});
    $.ajax({
        type: "POST",
        url: '/admin/transfer/transfer.html',
        data: $(".transfer-token-form").serializeArray(),
        success: function (data) {
            if (data.code >= 0) {
                //询问框
                alert("提交成功")

            } else if (data.code < 0) {
                alert(data.message);
            } else {
                alert("登录超时");
                window.location.href = "/public/login.html"
            }
            layer.close(index);

        }
    });
}

function transfer() {

    $("#transfer-modal").modal({
        relatedTarget: this,
        closeOnConfirm:false,
        onConfirm: function (e) {
            var coinType = e.data[0];
            var addresss = e.data[1];
            var fees = e.data[2];
            var password = e.data[3];
            var remark = e.data[4];

            var index = layer.load(0, {shade: false});
            $.ajax({
                type: "POST",
                url: '/admin/transfer/transfer.html',
                data: {
                    coinType: coinType,
                    address: addresss,
                    fees: fees,
                    password: password,
                    remark: remark
                },
                success: function (data) {
                    if (data.code >= 0) {
                        layer.msg("提交成功");
                    } else if (data.code < 0) {
                        layer.msg(data.message);
                    } else {
                        window.location.href = "/public/login.html"
                    }
                    layer.close(index);

                }
            });

        },
        onCancel: function (e) {

        }
    });

}

function initSetting() {
    var index = layer.load(0, {shade: false});
    $.ajax({
        type: "POST",
        url: '/public/initSetting.html',
        data: $(".init-setting-form").serializeArray(),
        success: function (data) {
            if (data.code >= 0) {

                // layer.open({
                //     type: 1,
                //     title:"请妥善保存以下信息",
                //     skin: 'layui-layer-rim', //加上边框
                //     area: ['500px', '500px'], //宽高
                //     content: '&nbsp;&nbsp;<strong>主账户地址：</strong><br/>' +data.data["address"]+
                //     '<br/>&nbsp;&nbsp;<strong>主账户KeyStore：</strong><br/>' + data.data["keystore"] + ''
                // });

                window.location.href = "/public/login.html"
            } else {
                alert(data.message);
            }

            layer.close(index);

        }
    });
}

function login() {
    var index = layer.load(0, {shade: false});

    $.ajax({
        type: "POST",
        url: '/public/login.html',
        data: $(".login-info-form").serializeArray(),
        success: function (data) {

            if (data.code >= 0) {
                window.location.href = "/admin/index.html"
            } else {
                alert(data.message);
            }

            layer.close(index);
        }
    });
}

var $addModal = $('#add-modal');
$addModal.siblings('.am-btn').on('click', function (e) {
    var $target = $(e.target);
    if (($target).hasClass('add-modal-open')) {
        $addModal.modal();
    }
});

var $editModal = $('#edit-modal');
$editModal.siblings('.am-btn').on('click', function (e) {
    var $target = $(e.target);
    if (($target).hasClass('edit-modal-open')) {
        $editModal.modal();
    }
});

function addManager() {
    var username = $("#doc-username-1").val();
    var password = $("#doc-password-1").val();

    $.ajax({
        type: "POST",
        url: '/admin/setting/addManager.html',
        data: {
            username: username,
            password: password,
            _xsrf: $("input[name=_xsrf]").val()
        },
        success: function (data) {
            alert(data.message);
        }
    });
}

function createAccounts() {

    $("#create-account-modal").modal({
        relatedTarget: this,
        onConfirm: function (e) {
            var coinType = e.data[0];
            var accountNumber = e.data[1];

            var index = layer.load(0, {shade: false});

            $.ajax({
                type: "POST",
                url: '/admin/account/createAccounts.html',
                data: {
                    accountNumber: accountNumber,
                    coinType: coinType
                },
                success: function (data) {

                    alert(data.message);

                    layer.close(index);

                }
            });

        },
        onCancel: function (e) {

        }
    });
}

function singleSummary() {

    $("#single-summary-modal").modal({
        relatedTarget: this,
        onConfirm: function (e) {
            var coinType = e.data[0];
            var accountNumber = e.data[1];

            var index = layer.load(0, {shade: false});

            $.ajax({
                type: "POST",
                url: '/admin/account/createAccounts.html',
                data: {
                    accountNumber: accountNumber,
                    coinType: coinType
                },
                success: function (data) {

                    alert(data.message);

                    layer.close(index);

                }
            });

        },
        onCancel: function (e) {

        }
    });
}

function forbidManager(id) {
    $.ajax({
        type: "POST",
        url: '/admin/setting/forbidManager.html',
        data: {
            id: id,
            _xsrf: $("input[name=_xsrf]").val()
        },
        success: function (data) {
            alert(data.message);
        }
    });
}

function updateSystemSetting() {
    var eth = $("#ethereum_rpc").val();
    var bit = $("#bitcoin_rpc").val();
    var lit = $("#litecoin_rpc").val();
    var time = $("#sync_balance_time").val();

    $.ajax({
        type: "POST",
        url: "/admin/setting/systemSetting.html",
        data: {
            ethereum_rpc: eth,
            bitcoin_rpc: bit,
            litcoin_rpc: lit,
            time: time,
            _xsrf: $("input[name=_xsrf]").val()
        },
        success: function (data) {
            alert(data.message)
        }
    });
}

String.prototype.endWith = function (s) {
    if (s == null || s == "" || this.length == 0 || s.length > this.length)
        return false;
    if (this.substring(this.length - s.length) == s)
        return true;
    else
        return false;
    return true;
}
String.prototype.startWith = function (s) {
    if (s == null || s == "" || this.length == 0 || s.length > this.length)
        return false;
    if (this.substr(0, s.length) == s)
        return true;
    else
        return false;
    return true;
}

$(".finish-transfer").each(function () {
    var clipboard = null;
    $(this).click(function () {
        clipboard = new Clipboard(".finish-transfer");//实例化
        clipboard.on('success', function (e) {
            layer.msg(e.text + " 已经复制到剪贴板");
            clipboard.destroy();
        });
    });
});