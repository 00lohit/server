import express from "express";
import { readCSV, writeCSV } from "./db.js";

const app = express();
const port = 3000;

app.use(express.json());

app.get("/items", async (req, res) => {
  try {
    const items = await readCSV();
    res.status(200).json({
      success: true,
      data: items,
      message: "Items retrieved successfully",
    });
  } catch (err) {
    console.error(err);
    res.status(500).json({ success: false, message: "Error reading data" });
  }
});

app.get("/items/:id", async (req, res) => {
  try {
    const items = await readCSV();
    const item = items.find((item) => item.id === req.params.id);
    if (item) {
      res.status(200).json({
        success: true,
        data: item,
        message: "Item retrieved successfully",
      });
    } else {
      res.status(404).json({ success: false, message: "Item not found" });
    }
  } catch (err) {
    console.error(err);
    res.status(500).json({ success: false, message: "Error reading data" });
  }
});

app.post("/items", async (req, res) => {
  try {
    const newItem = req.body;
    const items = await readCSV();

    let newId;
    if (items.length === 0) {
      newId = "1";
    } else {
      const lastItem = items[items.length - 1];
      newId = (parseInt(lastItem.id) + 1).toString();
    }

    newItem.id = newId;
    items.push(newItem);
    await writeCSV(items);
    res.status(201).json({
      success: true,
      data: newItem,
      message: "Item created successfully",
    });
  } catch (err) {
    console.error(err);
    res.status(500).json({ success: false, message: "Error creating item" });
  }
});

app.put("/items/:id", async (req, res) => {
  try {
    const updatedItem = req.body;
    const items = await readCSV();
    const index = items.findIndex((item) => item.id === req.params.id);

    if (index !== -1) {
      items[index] = { ...items[index], ...updatedItem };
      await writeCSV(items);
      res.status(200).json({
        success: true,
        data: items[index],
        message: "Item updated successfully",
      });
    } else {
      res.status(404).json({ success: false, message: "Item not found" });
    }
  } catch (err) {
    console.error(err);
    res.status(500).json({ success: false, message: "Error updating item" });
  }
});

app.delete("/items/:id", async (req, res) => {
  try {
    const items = await readCSV();
    const filteredItems = items.filter((item) => item.id !== req.params.id);

    if (filteredItems.length < items.length) {
      await writeCSV(filteredItems);
      res.status(200).json({
        success: true,
        message: "Item deleted successfully",
      });
    } else {
      res.status(404).json({ success: false, message: "Item not found" });
    }
  } catch (err) {
    console.error(err);
    res.status(500).json({ success: false, message: "Error deleting item" });
  }
});

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
