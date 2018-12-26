const api = "https://api.sampalm.com/";
sessid = undefined;
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
            AddtoChat("Error: "+res.Description);
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
            AddtoChat("Error: "+res.Description);
            return
        }
        AddtoChat(res.Description);
        sessid = res.Sessid; 
        AddtoChat("You are online now!")

        document.getElementById("loggin").style.display = "none";
        document.getElementById("logged").style.display = "block";

        ReadChat();
    });
}

function Logout(){
    let uname = document.getElementById("lg-username").value;
    qfetch("ch_logout", {Username: uname, Sessid: sessid}, "DELETE")
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("Error: "+res.Description);
            console.log(JSON.stringify(res))
            return
        }
        AddtoChat(res.Description);
        sessid = undefined; 
        AddtoChat("You are logged out now!")

        document.getElementById("loggin").style.display = "block";
        document.getElementById("logged").style.display = "none";
    });
}

function ReadChat(){
    if (sessid === undefined){
        return
    }
    let dt = {Sessid: sessid};
    if (lastchat !== undefined) {
        dt.LastID = lastchat.DateID
        dt.LastTime = lastchat.Time
    }
    console.log(JSON.stringify(dt))
    qfetch("ch_get_message", dt)
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("Error: "+res.Description);
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
    if(sessid === undefined){
        AddtoChat("Not Logged In")
        return
    }
    let txt = document.getElementById("ch-text").value;
    if (txt === ""){
        return;
    }

    qfetch("ch_message", {Sessid: sessid, Text: txt})
    .then(function(res){
        if(res.Value != 200){
            AddtoChat("Error: "+res.Description);
            return
        }
    })
}

// Update chat message
setInterval(function(){
    if (sessid !== undefined){
        qfetch("ch_get_message", {Sessid: sessid})
        .then(function(res){
            if(res.Value != 200){
                console.log("Error: "+res.Description);
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
for (let b in buttons) { 
    buttons[b].onclick = function(e){
        e.preventDefault;
        switch(buttons[b].name){
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
                Say();
                break;
            default:
                return;
        }    
    }
}
