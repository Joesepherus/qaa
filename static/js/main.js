function openCreateQuestionModal() {
  const modal = document.getElementById("create-question-modal");
  modal.classList.add("show");
}

function closeCreateQuestionModal() {
  const modal = document.getElementById("create-question-modal");
  modal.classList.remove("show");
}

function openCreateTrainingModal() {
  const modal = document.getElementById("create-training-modal");
  modal.classList.add("show");
}

function closeCreateTrainingModal() {
  const modal = document.getElementById("create-training-modal");
  modal.classList.remove("show");
}

function selectTrainingForRandomQuestion() {
  const trainingId = document.getElementById('training-select').value;
  console.log("trainingId", trainingId);
  window.location.href = "/random/" + trainingId
}

function openEditQuestionModal() {
  const modal = document.getElementById("edit-question-modal");
  modal.classList.add("show");
}

function closeEditQuestionModal() {
  const modal = document.getElementById("edit-question-modal");
  modal.classList.remove("show");
}

function selectQuestion(questionId) {
  window.location.href = "/questions/" + questionId
}

function isQuestionUrl() {
  const pathname = window.location.pathname;
  const regex = /^\/questions\/\d+$/;
  return regex.test(pathname);
}

function openEditTrainingModal() {
  const modal = document.getElementById("edit-training-modal");
  modal.classList.add("show");
}

function closeEditTrainingModal() {
  const modal = document.getElementById("edit-training-modal");
  modal.classList.remove("show");
}

function selectTraining(trainingId) {
  window.location.href = "/trainings/" + trainingId
}

function isTrainingURL() {
  const pathname = window.location.pathname;
  const regex = /^\/trainings\/\d+$/;
  return regex.test(pathname);
}

window.onload = function() {
  if (isQuestionUrl()) {
    openEditQuestionModal()
  }

  if (isTrainingURL()) {
    openEditTrainingModal()
  }
}
