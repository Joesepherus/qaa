function showModal() {
    const modal = document.getElementById("signup-modal");
    modal.classList.add("show");
    window.onclick = function(event) {
        if (event.target.classList.contains("modal")) {
            closeSignUpForm();
        }
    };
}

function closeSignUpForm() {
    const modal = document.getElementById("signup-modal");
    modal.classList.remove("show");
}

function showSignUpForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");
    signupForm.classList.add("active-form");
    loginForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showLoginForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");
    loginForm.classList.add("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showResetPasswordForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");

    loginForm.classList.remove("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.add("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showSetPasswordForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");

    loginForm.classList.remove("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.add("active-form");
}

function openModalShowLoginForm() {
    showModal();
    showLoginForm();
}

function openModalShowSignUpForm() {
    showModal();
    showSignUpForm();
}

function openModalShowSetPasswordForm() {
    showModal();
    showSetPasswordForm();
}

