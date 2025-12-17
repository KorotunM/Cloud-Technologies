function openMainModal() {
    document.getElementById('main-modal').style.display = 'block';
}

function openRegistrationModal() {
    closeModal('main-modal');
    document.getElementById('registration-modal').style.display = 'block';
}

function openLoginModal() {
    closeModal('main-modal');
    document.getElementById('login-modal').style.display = 'block';
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

window.onclick = function(event) {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
};