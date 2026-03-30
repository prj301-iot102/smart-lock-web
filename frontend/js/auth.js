window.addEventListener('DOMContentLoaded', () => {
  const token = localStorage.getItem('token');
  
  if (!token) {
    window.location.href = 'index.html';
  }
});