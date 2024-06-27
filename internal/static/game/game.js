let userScores = {}; // object to store scores for all users

ws.onopen = function() {
    console.log("Connected to WebSocket server");
};

ws.onmessage = function(event) {
    let data = JSON.parse(event.data);

    // update the specific user's score in the userScores object
    let user = data.ps;
    if (!userScores[user.id]) {
        userScores[user.id] = { fire: 0, water: 0, food: 0 };
        // create a new p element for the new user
        let userScore = document.createElement("p");
        userScore.id = `score-${user.id}`;
        document.getElementById("user-scores").appendChild(userScore);
    }
    userScores[user.id].fire = user.fire;
    userScores[user.id].water = user.water;
    userScores[user.id].food = user.food;

    // update the user's score display
    let userScore = document.getElementById(`score-${user.id}`);
    userScore.innerText = `${user.id}: fire: ${userScores[user.id].fire}, water: ${userScores[user.id].water}, food: ${userScores[user.id].food}`;

    // update the current user's display if it's the current user
    if (user.id === currentUserId) {
        document.getElementById("user-fire").innerText = "User Click: " + user.fire;
        document.getElementById("user-water").innerText = "User Click: " + user.water;
        document.getElementById("user-food").innerText = "User Click: " + user.food;
    }

    // find the maximum values among all users
    let maxValues = Object.values(userScores).reduce((max, user) => {
        return {
            fire: Math.max(max.fire, user.fire),
            water: Math.max(max.water, user.water),
            food: Math.max(max.food, user.food)
        };
    }, { fire: 0, water: 0, food: 0 });

    // update the maximum values display
    document.getElementById("max-fire").innerText = "Max Click: " + maxValues.fire;
    document.getElementById("max-water").innerText = "Max Click: " + maxValues.water;
    document.getElementById("max-food").innerText = "Max Click: " + maxValues.food;

    // update the global score display
    document.getElementById("fire").innerText = "Global Click: " + data.gs.fire;
    document.getElementById("water").innerText = "Global Click: " + data.gs.water;
    document.getElementById("food").innerText = "Global Click: " + data.gs.food;

    // explicitly add or remove the highlight class based on the maximum values
    const fireImg = document.getElementById("fire-item").querySelector("img");
    const waterImg = document.getElementById("water-item").querySelector("img");
    const foodImg = document.getElementById("food-item").querySelector("img");

    if (currentUserId in userScores) {
        const currentUser = userScores[currentUserId];

        if (currentUser.fire === maxValues.fire && currentUser.fire !== 0) {
            fireImg.classList.add("highlight");
        } else {
            fireImg.classList.remove("highlight");
        }

        if (currentUser.water === maxValues.water && currentUser.water !== 0) {
            waterImg.classList.add("highlight");
        } else {
            waterImg.classList.remove("highlight");
        }

        if (currentUser.food === maxValues.food && currentUser.food !== 0) {
            foodImg.classList.add("highlight");
        } else {
            foodImg.classList.remove("highlight");
        }
    }
};

ws.onclose = function() {
    console.log("Disconnected from WebSocket server");
};

ws.onerror = function(error) {
    console.error("WebSocket error:", error);
};

function clickItem(item) {
    // inject id for broadcast
    let msg = { "userId": currentUserId };
    msg[item] = 1;
    ws.send(JSON.stringify(msg));
}
