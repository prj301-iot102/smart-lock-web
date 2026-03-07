const loginForm = document.getElementById("loginForm");

loginForm.addEventListener("submit", async function(e){

    e.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try{
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/auth/login",{
            method:"POST",
            headers:{
                "Accept": "application/json, application/xml",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });

        console.log(response);

        const data = await response.json();

        console.log(data);

        if(response.ok) {
            localStorage.setItem("token", data.token);
            window.location.href="dashboard.html";
        }else{
            document.getElementById("error-msg").innerText = data.message || "Login failed";
        }
    }catch(error){
        console.log("Error:", error)
    }
});