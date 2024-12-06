import numpy as np
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation

# Simulate stock data for 20 companies over 4 years
np.random.seed(42)  # For reproducibility
companies = [f"Company {i+1}" for i in range(20)]
years = np.array([2020, 2021, 2022, 2023])

# Generate random stock prices in the range of 50 to 500 for 2020, with variability over the years
stock_prices = np.random.uniform(50, 500, (20, 4)) + np.random.normal(0, 20, (20, 4))

# Simulate stock changes due to military events (adding a hypothetical effect)
event_impact = np.array([0, -30, 50, -20])  # Impact per year (e.g., 2021 dip, 2022 growth, etc.)
stock_prices += event_impact

# Normalize for better visualization
normalized_prices = (stock_prices - np.min(stock_prices)) / (np.max(stock_prices) - np.min(stock_prices))

# Set up the plot
fig, ax = plt.subplots(figsize=(14, 8))
lines = []
for _ in companies:
    line = ax.plot([], [], alpha=0.7)[0]
    lines.append(line)

# Highlight the years of significant military events
for year, impact in zip(years, event_impact):
    ax.axvline(x=year, color='gray', linestyle='--', alpha=0.5)
    ax.text(year, 1.05, f"{impact:+} Impact", ha='center', fontsize=10, color='red')

ax.set_xlim(years[0], years[-1])
ax.set_ylim(0, 1.1)
ax.set_title("Simulated Stock Changes of Top 20 Fortune 500 Companies (2020-2023)")
ax.set_xlabel("Year")
ax.set_ylabel("Normalized Stock Prices")
ax.grid(alpha=0.3)

# Initialize function
def init():
    for line in lines:
        line.set_data([], [])
    return lines

# Animation function
def animate(frame):
    year_idx = frame
    for i, line in enumerate(lines):
        line.set_data(years[:year_idx + 1], normalized_prices[i, :year_idx + 1])
    return lines

# Create the animation
ani = FuncAnimation(fig, animate, frames=len(years), init_func=init, blit=True, interval=1000, repeat=False)

# Show the animation
plt.show()
