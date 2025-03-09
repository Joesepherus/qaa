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

function selectTraining() {
  const trainingId = document.getElementById('training-select').value;
  window.location.href = "/random/" + trainingId
}

