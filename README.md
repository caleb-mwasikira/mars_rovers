In the vast expanses of the Martian terrain, multiple rovers have been deployed for scientific exploration. These 
rovers traverse the Martian landscape, collecting valuable data for research purposes. Understanding the spatial 
relationships between these rovers is crucial for optimizing their operations and avoiding potential collisions.

**Problem:**
Using the programming language of your choice, develop an algorithm that calculates the nearest and farthest rovers 
within a given dataset. The dataset contains information about the position coordinates (latitude and longitude) of 
each rover at specific time intervals during their mission on Mars.

**Dataset Format:**
Mars rovers dataset
SHA256sum:
00914983a3f89cde7b538f2d7cd54675022637324b6221f8d71c6714a16883ce  mars_rovers.csv

The dataset is structured as follows:
- Rover ID: Unique identifier for each rover.
- Latitude: Geographical coordinate indicating the north-south position of the rover.
- Longitude: Geographical coordinate indicating the east-west position of the rover.
- Timestamp: Time at which the rover's position was recorded.

--------

Planet size dataset
SHA256sum: 6dde0b908655aafa643d8c1010d0c21e6ceda796a1ddd6080867b54f8bc8dd14  planets_size.csv

The dataset is structured as follows:
- Planet: Name of the planet
- Radius: Radius of the planet in kilometres

**Constraints and Considerations:**
- The dataset may contain irregular time intervals between position recordings.
- Rovers may have different trajectories and may not necessarily cover the same geographical areas.
- The solution should be scalable for datasets with a varying number of rovers and time points.
- Ensure the efficiency and computational feasibility of the algorithm, especially for large datasets.

**Expected Output:**
The output should include the Rover IDs and the corresponding distance for both the nearest and farthest rover pairs. 
The solution should be capable of handling diverse datasets representing the dynamic movements of rovers on the 
Martian surface.
