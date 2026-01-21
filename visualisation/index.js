const REFRESH_INTERVAL_MS = 500; // частота опроса

async function updateLoop() {
  try {
    gripper = await fetchGripper()
    carousel = await fetchCarousel()
    start = await fetchStart()
    packaging = await fetchPackaging()
    sorting = await fetchSorting()

    gripper_puck = document.getElementById("gripper-puck")
    gripper_hor_pos = document.getElementById("gripper-hor-pos")
    gripper_ver_pos = document.getElementById("gripper-ver-pos")

    start_slot = document.getElementById("start-slot")
    package_slot = document.getElementById("package-slot")
    sort_slot = document.getElementById("sort-slot")

    carousel_slot1 = document.getElementById("car-slot1")
    carousel_slot2 = document.getElementById("car-slot2")
    carousel_slot3 = document.getElementById("car-slot3")
    carousel_slot4 = document.getElementById("car-slot4")
    carousel_slot5 = document.getElementById("car-slot5")
    carousel_slot6 = document.getElementById("car-slot6")

    carousel_slots = [carousel_slot1, carousel_slot2, carousel_slot3, carousel_slot4, carousel_slot5, carousel_slot6]

    // carousel
    arr = carousel.slots
    for (let i = 0; i < arr.length; i++) {
        if (arr[i] == null) {
            carousel_slots[i].style.background = "burlywood"
        } else {
            carousel_slots[i].style.background = "black"
        }
    }

    // gripper
    if (gripper.puckSlot == null) {
        gripper_puck.style.visibility = "hidden"
    } else {
        gripper_puck.style.visibility = "visible"
    }
    gripper_hor_pos.innerHTML = Number((gripper.curHorizontalPosition).toFixed(2)); // 6.7
    gripper_ver_pos.innerHTML = Number((gripper.curVerticalPosition).toFixed(2)); // 6.7

    // start
    if (start.puckSlot == null) {
        start_slot.style.background = "burlywood"
    } else {
        start_slot.style.background = "black"
    }

    // package
    if (packaging.puckSlot == null) {
        package_slot.style.background = "burlywood"
    } else {
        package_slot.style.background = "black"
    }

    // sorting
    if (sorting.puckSlot == null) {
        sort_slot.style.background = "burlywood"
    } else {
        sort_slot.style.background = "black"
    }

  } catch (e) {
    console.error("Failed to update state", e);
  }
}

async function fetchGripper() {
    const response = await fetch(`http://localhost:8080/vis/gripper`);
    if (!response.ok) {
        console.log(response.status, response.statusText)
    }
    return response.json();
}

async function fetchCarousel() {
    const response = await fetch(`http://localhost:8080/vis/carousel`);
    if (!response.ok) {
        console.log(response.status, response.statusText)
    }
    return response.json();
}

async function fetchStart() {
    const response = await fetch(`http://localhost:8080/vis/start`);
    if (!response.ok) {
        console.log(response.status, response.statusText)
    }
    return response.json();
}

async function fetchPackaging() {
    const response = await fetch(`http://localhost:8080/vis/packaging`);
    if (!response.ok) {
        console.log(response.status, response.statusText)
    }
    return response.json();
}

async function fetchSorting() {
    const response = await fetch(`http://localhost:8080/vis/sorting`);
    if (!response.ok) {
        console.log(response.status, response.statusText)
    }
    return response.json();
}

setInterval(updateLoop, REFRESH_INTERVAL_MS);
updateLoop()