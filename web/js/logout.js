document.addEventListener("DOMContentLoaded", function () {
    const logoutBtn = document.getElementById("logout-btn");

    if (logoutBtn) {
        logoutBtn.addEventListener("click", function () {
            fetch("/logout", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                }
            })
            .then(response => {
                if (response.ok) {
                    updateUserUI(null);
                } else {
                    console.error("Ошибка при выходе из системы");
                }
            })
            .catch(err => {
                console.error("Ошибка при запросе на выход:", err);
            });
        });
    }
});

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        return parts.pop().split(';').shift();
    }
    return null;
}

// Функция для обновления UI
function updateUserUI(email) {
    const loginBtn = document.getElementById("login-btn");
    const userInfoDiv = document.getElementById("user-info");
    const userEmailSpan = document.getElementById("user-email");

    if (email) {
        userEmailSpan.textContent = email;
        userInfoDiv.style.display = "flex";
        loginBtn.style.display = "none";
    } else {
        userEmailSpan.textContent = "";
        userInfoDiv.style.display = "none";
        loginBtn.style.display = "block";
    }
}

// Вызов при загрузке страницы
document.addEventListener("DOMContentLoaded", function () {
    const userEmail = getCookie("user_email");
    updateUserUI(userEmail);
});
