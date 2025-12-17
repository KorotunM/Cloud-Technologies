document.addEventListener("DOMContentLoaded", function () {
    function getQueryParams() {
        const params = new URLSearchParams(window.location.search);
        return {
            email: params.get("email"),
            status: params.get("status"),
            error: params.get("error"),
        };
    }

    function handleOAuthResponse() {
        const { email, status, error } = getQueryParams();

        if (status === "success" && email) {
            updateUserUI(email);
            history.replaceState(null, "", "/"); // Очистим URL от query параметров
        } else if (error) {
            alert("Ошибка: " + error);
        }
    }

    handleOAuthResponse();
});

