const api = "https://api.sampalm.com/";
let lastchat = undefined;

function qfetch(lambdaName, data, method = "POST") {
    return fetch(api + lambdaName, {
        headers: { "content-type": "application/json" },
        method: method,
        body: JSON.stringify(data)
    }).then(function (res) {
        return res.json();
    });
}

function AddtoChat(txt) {
    let ctdiv = document.getElementById("chat-box");
    let p = document.createElement("p");
    p.innerHTML = txt;
    ctdiv.appendChild(p);
    ctdiv.scrollTop = ctdiv.offsetHeight * 100;
}

function NewUser() {
    let uname = document.getElementById("lg-username").value;
    let pass = document.getElementById("lg-pass").value;
    if (uname.length === 0 || pass.length === 0) {errorToHtml("chat", "Invalid user or password", "error"); return;}

    qfetch("ch_new_user", { Username: uname, Password: pass })
        .then(function (res) {
            if (res.Value != 200) {
                errorToHtml("chat", res.Description, "error")
                return
            }
            errorToHtml("chat", res.Description, "info")
        });
}

function Login() {
    let uname = document.getElementById("lg-username").value;
    let pass = document.getElementById("lg-pass").value;
    if (uname.length === 0 || pass.length === 0) {errorToHtml("chat", "Invalid user or password", "error"); return;}

    qfetch("ch_login", { Username: uname, Password: pass })
        .then(function (res) {
            if (res.Value != 200) {
                errorToHtml("chat", res.Description, "error")
                return
            }
            errorToHtml("chat", res.Description, "info")
            document.cookie = "sessid=" + res.Sessid;
            toogleBtn();
            ReadChat();
        });
}

function Logout() {
    let uname = document.getElementById("lg-username").value;
    errorToHtml("chat", "You are being logged out", "info");
    qfetch("ch_logout", { Username: uname, Sessid: getCookie() }, "DELETE")
        .then(function (res) {
            document.cookie = "sessid=; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
            toogleBtn(false);
        });
}

function ReadChat() {
    if (getCookie() === undefined) {
        return
    }
    let dt = { Sessid: getCookie() };
    if (lastchat !== undefined) {
        dt.LastID = lastchat.DateID
        dt.LastTime = lastchat.Time
    }
    console.log(JSON.stringify(dt))
    qfetch("ch_get_message", dt)
        .then(function (res) {
            if (res.Value != 200) {
                AddtoChat("<span class='msg-error'>Error: " + res.Description + "</span>");
                return
            }
            if (res.Chats !== undefined) {
                res.Chats.map(chat => {
                    AddtoChat(`<b>${chat.Username}</b> says:\n
                    <span class="chat-txt">${chat.Text}</span><i class="chat-time">${chat.Time}</i>`);
                    lastchat = chat;
                })
            }
        })
}

function Say() {
    if (getCookie() === undefined) {
        AddtoChat("<span class='msg-error'>Not Logged In!</span>")
        return
    }
    let txt = document.getElementById("ch-text").value;
    if (txt === "") {
        return;
    }

    qfetch("ch_message", { Sessid: getCookie(), Text: txt })
        .then(function (res) {
            if (res.Value != 200) {
                AddtoChat("<span class='msg-error'>Error: " + res.Description + "</span>");
                return
            }
        })
}

function Translate() {
    let txt = document.getElementById("ch-text");
    let src = document.getElementById("ch-src").value;
    let tg = document.getElementById("ch-tg").value;
    if (txt === "") {
        return;
    }
    if (getCookie() === undefined) {
        AddtoChat("<span class='msg-error'>Not Logged In!</span>")
        return;
    }
    qfetch("ch_translate", { Sessid: getCookie(), Source: src, Target: tg, Text: txt.value })
        .then(function (res) {
            if (res.Value != 200) {
                AddtoChat("<span class='msg-error'>Error: " + res.Description + "</span>");
                return
            }
            txt.value = res.Body;
        });
}

function toogleBtn(on = true) {
    if (on) {
        AddtoChat("You are online now!")

        document.getElementById("loggin").style.display = "none";
        document.getElementById("logged").style.display = "block";
        return;
    }
    AddtoChat("You are logged out!");

    document.getElementById("loggin").style.display = "block";
    document.getElementById("logged").style.display = "none";
}

function getCookie() {
    return decodeURIComponent(document.cookie).split("sessid=")[1] !== undefined ?
        decodeURIComponent(document.cookie).split("sessid=")[1].substring(0, 40) :
        undefined;
}

function errorToHtml(id, error, type, persist=false) {
    removeErrors();
    let selector = type === "error" ? 
    `#${id}` : type === "alert" ?
     `#${id} .col-12` : type === "info" ?
     `#${id}` : `#${id}`;
    let div = document.querySelector(selector);
    let p = document.createElement("p");
    p.className = "warning warning-" + type;
    p.innerHTML = error;
    div.insertAdjacentElement("afterbegin", p)
    if(persist) return;
    setTimeout(function(){
        p.classList.add("fade");
        setTimeout(p.remove(), 500)
    }, 2000);
}
function removeErrors() {
    let old = document.getElementsByClassName("warning");
    for (o of old) {
        o.parentNode.removeChild(o);
    }
}
function removeTablesResults() {
    let tbd = document.querySelector("#files table tbody");
    tbd.innerHTML = "";
}

function listBucket() {
    let table = document.querySelector("#files table");
    if (getCookie() === undefined) {
        table.style.display = "none";
        removeTablesResults();
        errorToHtml("files", "You must be <a href='#chat'>logged in</a> to see your files.", "alert", true);
        return
    }
    removeErrors();
    qfetch("ch_list_bucket", { Sessid: getCookie() })
        .then(function (res) {
            if (res.Value !== 200) {
                errorToHtml("files", res.Description, "error");
                return
            }
            if (res.Body.length === 0) {
                table.style.display = "none";
                removeTablesResults();
                errorToHtml("files", "You haven't any file stored yet.", "alert");
                return;
            }
            table.style.display = "table";
            let tb = document.getElementById("files-list");
            tb.innerHTML = "";
            res.Body.map(file => {
                let tr = document.createElement("tr");
                let ln = `<td>${file.Key}</td>
                <td class="center-text">${file.Size}KB</td>
                <td class="center-text file-download"><i class="fas fa-cloud-download-alt" data-foo="${file.Key}"></i>
                ${file.Key.split(".").pop() === "mp3" ? `<i class="fas fa-play-circle play" id="${file.Key.hashCode()}" data-foo="${file.Key}"></i></td>` : '</td>'}`;
                tr.innerHTML = ln;
                tb.appendChild(tr);
            });
            document.body.addEventListener("click", function loader(e) {
                if (e.srcElement.hasAttribute("data-foo")) {
                    let filename = e.srcElement.getAttribute("data-foo");
                    e.srcElement.removeAttribute("data-foo");
                    if (e.srcElement.classList.contains("play")) {
                        getObject(filename, "play");
                        e.srcElement.setAttribute("data-foo", filename)
                        return;
                    }
                    getObject(filename);
                    setTimeout(() => (e.srcElement.setAttribute("data-foo", filename)), 5000)
                    return;
                }
            }, false);
        })
}


String.prototype.hashCode = function() {
    var hash = 0, i, chr;
    if (this.length === 0) return hash;
    for (i = 0; i < this.length; i++) {
      chr   = this.charCodeAt(i);
      hash  = ((hash << 5) - hash) + chr;
      hash |= 0;
    }
    return hash;
};
function base64ToArrayBuffer(data) {
    var binaryString = window.atob(data);
    var binaryLen = binaryString.length;
    var bytes = new Uint8Array(binaryLen);
    for (var i = 0; i < binaryLen; i++) {
        var ascii = binaryString.charCodeAt(i);
        bytes[i] = ascii;
    }
    return bytes;
}

let lastAudioID = undefined;
function getObject(file, action="download") {
    if (!getCookie()) {
        return;
    }
    qfetch("ch_get_obj", { SessID: getCookie(), Filename: file })
        .then(function (res) {
            if (res.Value !== 200) {
                console.log(res.Description)
                return;
            }
            var arrBuffer = base64ToArrayBuffer(res.Body);
            let blob = new Blob([arrBuffer], { type: res.ContentType });
            let filename = file.split('\\').pop().split('/').pop();
            if (action === "play") {
                let fid = filename.hashCode();
                let audioplayer = document.getElementById("="+fid);
                let audiobtn = document.getElementById(fid);
                if(audioplayer){
                    if(audioplayer.paused){
                        audiobtn.classList.remove("fa-play-circle");
                        audiobtn.classList.add("fa-pause-circle");
                        audioplayer.play()
                        return;
                    }
                    audiobtn.classList.add("fa-play-circle");
                    audiobtn.classList.remove("fa-pause-circle");
                    audioplayer.pause();
                    return;
                }
                if (lastAudioID !== undefined) {
                    audiobtn = document.getElementById(lastAudioID);
                    audiobtn.classList.add("fa-play-circle");
                    audiobtn.classList.remove("fa-pause-circle");
                    lastAudioID = undefined;
                }
                playFile(blob, filename)
                return;
            }
            saveFile(blob, filename);
            return; 
    });
}

function saveFile(blob, filename) {
    if (window.navigator.msSaveOrOpenBlob) {
        window.navigator.msSaveOrOpenBlob(blob, filename);
    } else {
        const a = document.createElement('a');
        document.body.appendChild(a);
        const url = window.URL.createObjectURL(blob);
        a.href = url;
        a.download = filename;
        a.click();
        setTimeout(() => {
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        }, 0)
    }
}

function playFile(blob, filename) {
    let audiobtn = document.getElementById(filename.hashCode());
    const a = document.querySelector("audio") || document.createElement('audio');
    document.body.appendChild(a);
    const url = window.URL.createObjectURL(blob);
    a.type = blob.type;
    a.src = url;
    a.id = "="+filename.hashCode();
    a.play();
    a.onended = function() {
        audiobtn.classList.add("fa-play-circle");
        audiobtn.classList.remove("fa-pause-circle");
    }
    audiobtn.classList.add("fa-pause-circle");
    lastAudioID = filename.hashCode();
}

(function () {
    if (getCookie() !== undefined) {
        toogleBtn();
        ReadChat();
        return;
    }
    toogleBtn(false);
})();

setInterval(function () {
    // Update chat messages
    if (getCookie() !== undefined) {
        qfetch("ch_get_message", { Sessid: getCookie() })
            .then(function (res) {
                if (res.Value !== 200) {
                    console.log(res.Description);
                    return
                }
                if (res.Chats !== undefined) {
                    let newchat = res.Chats[res.Chats.length - 1]
                    if (lastchat === undefined) {
                        return;
                    }
                    console.log("NewTime: " + newchat.Time + " - LastTime: " + lastchat.Time)
                    if (newchat.Time !== lastchat.Time) {
                        console.log("Calling ReadChat()");
                        ReadChat();
                    }
                }
            })
    }
}, 1500);
let buttons = document.querySelectorAll("input[type=button]");
for (let b of buttons) {
    b.onclick = function (e) {
        e.preventDefault;
        switch (b.name) {
            case ("lg-user"):
                this.disabled = true;
                errorToHtml("chat", "Log in to your account", "info")
                Login();
                setTimeout(() => this.disabled = false, 2000);
                break;
            case ("lg-out"):
                Logout();
                break;
            case ("new-user"):
                this.disabled = true;
                errorToHtml("chat", "Creating new user", "info")
                NewUser();
                setTimeout(() => this.disabled = false, 2000);
                break;
            case ("ch-send"):
                this._lockButton("Sending ");
                Say();
                setTimeout(() => this._unlockButton("Send Message"), 2000);
                break;
            case ("ch-tl"):
                this._lockButton("Translating ");
                Translate();
                setTimeout(() => this._unlockButton("Translate Message"), 2000);
                break;
            default:
                return;
        }
    }
}
Element.prototype._lockButton = function(text){
    this.disabled = true;
    this.value = text || this.value;
    this.classList.add("loading", "btn-pd");
}
Element.prototype._unlockButton = function(text){
    this.disabled = false; 
    this.value = text || this.value;; 
    this.classList.remove("btn-pd", "loading");
}
let fls = document.getElementById("files");
fls.onload = listBucket();