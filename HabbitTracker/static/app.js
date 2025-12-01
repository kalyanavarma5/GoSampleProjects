const habitForm = document.getElementById('habitForm');
const habitList = document.getElementById('habitList');

function fetchHabits() {
    fetch('/habits')
        .then(res => {
            if (!res.ok) {
                throw new Error(`HTTP ${res.status}: ${res.statusText}`);
            }
            return res.json();
        })
        .then(habits => {
            habitList.innerHTML = '';
            habits.forEach(habit => {
                const li = document.createElement('li');

                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.checked = habit.completed;
                checkbox.onchange = () => {
                    fetch(`/toggle?id=${habit.id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' }
                    })
                        .then(res => {
                            if (!res.ok) {
                                throw new Error(`Toggle failed: ${res.status}`);
                            }
                            fetchHabits(); // Refresh list
                        })
                        .catch(err => console.error('Toggle error:', err));
                };

                li.appendChild(checkbox);
                li.appendChild(document.createTextNode(' ' + habit.name));
                habitList.appendChild(li);
            });
        })
        .catch(err => {
            console.error('Fetch habits error:', err);
            habitList.innerHTML = '<li style="color:red">Error loading habits</li>';
        });
}


habitForm.onsubmit = (e) => {
    e.preventDefault();
    const habitName = document.getElementById('habitName').value;
    fetch('/habits', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: habitName, completed: false }),
    })
        .then(res => res.json())
        .then(() => {
            habitForm.reset();
            fetchHabits();
        });
};

// Initial fetch when page loads
fetchHabits();
