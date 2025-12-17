
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

function registerUser() {
    const email = document.getElementById("reg-email").value.trim();
    const password = document.getElementById("reg-password").value.trim();
    const confirmPassword = document.getElementById("confirm-password").value.trim();

    if (!email || !password || !confirmPassword) {
        alert("Все поля обязательны для заполнения.");
        return;
    }

    if (password !== confirmPassword) {
        alert("Пароли не совпадают.");
        return;
    }

    fetch("/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
    })
        .then((res) => res.json())
        .then((data) => {
            if (data.status === "success") {
                updateUserUI(data.email);
                closeModal('registration-modal');
            } else {
                alert(data.error);
            }
        })
        .catch((err) => {
            console.error("Ошибка при регистрации:", err);
        });
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}