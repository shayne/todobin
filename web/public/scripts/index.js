function getStoredLists() {
  var lists = window.localStorage.getItem('todo_lists');

  return lists ? JSON.parse(lists) : {}
}

function saveLists(lists) {
  window.localStorage.setItem('todo_lists', JSON.stringify(lists))
}
