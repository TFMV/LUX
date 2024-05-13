# LUX - Lightweight Update eXchange

## Overview
LUX is a real-time data analytics platform that leverages the power of Go and Apache Iceberg to process IoT data efficiently. It utilizes Google Cloud Pub/Sub for messaging and Iceberg tables hosted on Google Cloud Storage (GCS) to ensure scalability and reliability in data handling.

![LUX Architecture](lux.webp)

## Key Features

- **Real-Time Data Processing:** Stream data from IoT devices directly into Iceberg tables for immediate analysis.
- **Scalability:** Easily scale your data processing workload with the elasticity of Google Cloud services.
- **Flexibility:** Deployable on Google Cloud Run, providing a serverless environment that scales automatically.

## Getting Started

1. **Set up Google Cloud Services:**
   - Ensure that you have a Google Cloud project set up with Pub/Sub and Cloud Run enabled.
   - Configure your GCS bucket and Pub/Sub topic as described in the provided Terraform scripts.

2. **Deploy the Application:**
   - Clone the repository and navigate to the deployment directory.
   - Use the provided Dockerfile to build and deploy your application to Google Cloud Run.

3. **Monitor and Scale:**
   - Monitor your application's performance directly from Google Cloud Console.
   - Adjust the scaling settings in Cloud Run to meet your demand.

## Contributions

Contributions are welcome! Please feel free to submit pull requests, or file issues for bugs, questions, and feature requests.

## License

This project is licensed under the terms of the MIT license.

## Acknowledgments

Special thanks to the open-source community and everyone who has contributed to the development of Apache Iceberg and Google Cloud Pub/Sub.
