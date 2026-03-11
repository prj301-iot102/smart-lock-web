document.addEventListener("click", function(e) {
    if (e.target.id === "logoutBtn") {
        localStorage.removeItem("token");
        window.location.href = "./index.html";
    }
});