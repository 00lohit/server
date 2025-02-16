import fs from "fs";
import csv from "csv-parser";
import { createObjectCsvWriter } from "csv-writer";

const csvFilePath = "data.csv";

async function readCSV() {
  return new Promise((resolve, reject) => {
    const results = [];
    fs.createReadStream(csvFilePath)
      .pipe(csv())
      .on("data", (data) => results.push(data))
      .on("end", () => resolve(results))
      .on("error", reject);
  });
}

async function writeCSV(data) {
  return new Promise((resolve, reject) => {
    if (!data || data.length === 0) {
      return resolve();
    }
    const headers = Object.keys(data[0]);
    const csvWriter = createObjectCsvWriter({
      path: csvFilePath,
      header: headers.map((header) => ({ id: header, title: header })),
    });

    csvWriter
      .writeRecords(data)
      .then(() => resolve())
      .catch(reject);
  });
}

export { readCSV, writeCSV };
