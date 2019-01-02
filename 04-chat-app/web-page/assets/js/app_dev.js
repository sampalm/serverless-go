const api = "https://api.sampalm.com/";
lastchat = undefined;

function qfetch(lambdaName, data, method="POST"){
    return fetch(api+lambdaName, {
        headers:{"content-type": "application/json"},
        method: method,
        body: JSON.stringify(data)
    }).then(function(res){
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

function NewUser(){
    let uname = document.getElementById("lg-username").value;
    AddtoChat("Creating User: "+uname);
    let pass = document.getElementById("lg-pass").value;

    qfetch("ch_new_user", {Username:uname, Password:pass})
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("<span class='msg-error'>Error: "+res.Description+"</span>");
            return
        }
        AddtoChat(res.Description);
    });
}

function Login(){
    let uname = document.getElementById("lg-username").value;
    AddtoChat("Log in User: "+uname);
    let pass = document.getElementById("lg-pass").value;

    qfetch("ch_login", {Username: uname, Password: pass})
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("<span class='msg-error'>Error: "+res.Description+"</span>");
            return
        }
        AddtoChat(res.Description);
        document.cookie = "sessid="+res.Sessid; 
        toogleBtn();
        ReadChat();
    });
}

function Logout(){
    let uname = document.getElementById("lg-username").value;
    qfetch("ch_logout", {Username: uname, Sessid: getCookie()}, "DELETE")
    .then(function(res){
        document.cookie = "sessid=; expires=Thu, 01 Jan 1970 00:00:00 UTC;"; 
        toogleBtn(false);
    });
}

function ReadChat(){
    if (getCookie() === undefined){
        return
    }
    let dt = {Sessid: getCookie()};
    if (lastchat !== undefined) {
        dt.LastID = lastchat.DateID
        dt.LastTime = lastchat.Time
    }
    console.log(JSON.stringify(dt))
    qfetch("ch_get_message", dt)
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("<span class='msg-error'>Error: "+res.Description+"</span>");
            return
        }
        if(res.Chats !== undefined){
            for(let c in res.Chats){
                let ch = res.Chats[c];
                AddtoChat("<b>"+ch.Username + "</b> says:");
                AddtoChat("<span class='chat-txt'>"+ch.Text+"</span>"+"<i class='chat-time'>"+ch.Time+"</i>");
                lastchat = ch;
            }
        }
    })
}

function Say(){
    if(getCookie() === undefined){
        AddtoChat("<span class='msg-error'>Not Logged In!</span>")
        return
    }
    let txt = document.getElementById("ch-text").value;
    if (txt === ""){
        return;
    }

    qfetch("ch_message", {Sessid: getCookie(), Text: txt})
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("<span class='msg-error'>Error: "+res.Description+"</span>");
            return
        }
    })
}

function Translate() {
    let txt = document.getElementById("ch-text");
    let src = document.getElementById("ch-src").value;
    let tg = document.getElementById("ch-tg").value;
    if (txt === ""){
        return;
    }
    if (getCookie() === undefined){
        AddtoChat("<span class='msg-error'>Not Logged In!</span>")
        return;
    }
    console.log("translating...")
    qfetch("ch_translate", {Sessid: getCookie(), Source: src, Target: tg, Text: txt.value})
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("<span class='msg-error'>Error: "+res.Description+"</span>");
            return
        }
        txt.value = res.Body;
    });
}

function toogleBtn(on=true){
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

function getCookie(){
    return decodeURIComponent(document.cookie).split("sessid=")[1] !== undefined ? 
    decodeURIComponent(document.cookie).split("sessid=")[1].substring(0,40) :
    undefined;
}

function listBucket(){
    if (getCookie() !== undefined){
        qfetch("ch_list_bucket", {Sessid: getCookie()})
        .then(function(res){
            if(res.Value != 200){
                console.log(res.Description)
                return
            }
            let tb = document.getElementById("files-list");
            tb.innerHTML = "";
            for(let file of res.Body){
                let tr = document.createElement("tr");
                let ln = `<td>${file.Key}</td>`;
                ln += `<td>${file.Size}</td>`;
                ln += `<td class="file-download"><i class="fas fa-cloud-download-alt" data-foo="${file.Key}"></i></td>`;
                tr.innerHTML = ln;
                tb.appendChild(tr);
            }  
            document.body.addEventListener("click", function loader(e){
                if(e.srcElement.hasAttribute("data-foo")){
                    let filename = e.srcElement.getAttribute("data-foo");
                    console.log("Downloading...")
                    e.srcElement.removeAttribute("data-foo");
                    getObject(filename);
                    setTimeout(()=>(e.srcElement.setAttribute("data-foo", filename)), 5000)
                }
            }, false);
        })
    }
}

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

function getObject(file){
    if(!getCookie()){
        return;
    }
    qfetch("ch_get_obj", {SessID: getCookie(), Filename: file})
    .then(function(res) {
        if(res.Value !== 200){
            console.log(res.Description)
            return
        }
        var arrBuffer = base64ToArrayBuffer(res.Body);
        let blob = new Blob([arrBuffer], {type: res.ContentType});
        let filename = file.split('\\').pop().split('/').pop();
        saveFile(blob, filename)
    })
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

(function(){
    console.log("CookieState: "+getCookie());
    if (getCookie() !== undefined){
        toogleBtn();
        ReadChat();
        return;
    }
    toogleBtn(false);
})();

// Update chat message
setInterval(function(){
    if (getCookie() !== undefined){
        qfetch("ch_get_message", {Sessid: getCookie()})
        .then(function(res){
            if(res.Value != 200){
                console.log(res.Description);
                return
            }
            if(res.Chats !== undefined){
                let newchat = res.Chats[res.Chats.length-1]
                if (lastchat === undefined){
                    return;
                }
                console.log("NewTime: "+newchat.Time +" - LastTime: "+lastchat.Time)
                if(newchat.Time !== lastchat.Time) {
                    console.log("Calling ReadChat()");
                    ReadChat();
                }
            }
        })
    }
}, 1500);
let buttons = document.querySelectorAll("input[type=button]"); 
for (let b of buttons) { 
    b.onclick = function(e){
        e.preventDefault;
        switch(b.name){
            case ("lg-user"):
                Login();
                break;
            case ("lg-out"):
                Logout();
                break;
            case ("new-user"):
                NewUser();
                break;
            case ("ch-send"):
                this.disabled = true;
                this.value = "Sending  "
                this.classList.add("loading", "btn-pd");
                Say();
                setTimeout(()=>{this.value = "Send Message"; this.disabled = false; this.classList.remove("btn-pd", "loading");}, 2000);
                break;
            case ("ch-tl"):
                this.disabled = true;
                this.value = "Translating  "
                this.classList.add("loading", "btn-pd");
                Translate();
                setTimeout(()=>{this.value = "Translate Message"; this.disabled = false; this.classList.remove("btn-pd", "loading");}, 2000);
                break;
            default:
                return;
        }    
    }
}