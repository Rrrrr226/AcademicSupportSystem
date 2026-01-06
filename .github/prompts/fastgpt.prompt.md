---
name: fastgpt-api
description: FastGPT API doc
agent: 'agent'
---

openapi: 3.0.0
info:
  title: FastGPT Open Source API
  description: >
    API documentation for FastGPT non-commercial features.
    
    **Base URL Note**: The BaseURL is the root address for all interfaces (e.g., `https://api.fastgpt.in/api`).
    
    **Authentication**:
    All interfaces require the `Authorization` header.
    - **App-specific Key**: Used for Chat Completion interfaces (starts with `fastgpt-`).
    - **Global Key**: Used for management interfaces (dataset, app management, etc.).
  version: 4.8.x
servers:
  - url: http://localhost:3000/api
    description: Local development environment
  - url: https://api.fastgpt.in/api
    description: Official SaaS (Example)
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: API Key
  schemas:
    ChatCompletionRequest:
      type: object
      required:
        - messages
      properties:
        chatId:
          type: string
          description: Chat ID. If provided, history is used/stored. If empty, it's a stateless conversation.
        stream:
          type: boolean
          default: false
          description: Whether to stream the response.
        detail:
          type: boolean
          default: false
          description: Whether to return detailed execution information (module status, references, etc.).
        variables:
          type: object
          description: Module variables to replace `[key]` in inputs.
        messages:
          type: array
          items:
            type: object
            properties:
              role:
                type: string
                enum: [user, assistant, system]
              content:
                type: string
                description: Text content. (For vision/files, use array structure as per OpenAI spec).
    
    ChatResponse:
      type: object
      properties:
        id:
          type: string
        choices:
          type: array
          items:
            type: object
            properties:
              message:
                type: object
                properties:
                  role:
                    type: string
                  content:
                    type: string
              finish_reason:
                type: string
    
    DatasetCreateRequest:
      type: object
      required:
        - name
      properties:
        parentId:
          type: string
          nullable: true
          description: Parent folder ID.
        type:
          type: string
          enum: [dataset, folder]
          default: dataset
        name:
          type: string
        intro:
          type: string
        avatar:
          type: string
        vectorModel:
          type: string
          description: Vector model name (leave empty for system default).
        agentModel:
          type: string
          description: QA generation model (leave empty for system default).

security:
  - BearerAuth: []

paths:
  # ============================
  # Chat Interface
  # ============================
  /v1/chat/completions:
    post:
      summary: Chat Completion (Standard & Workflow)
      description: >
        Compatible with OpenAI Chat Completion API.
        Requires **App-specific Key**.
      tags:
        - Chat
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ChatCompletionRequest'
      responses:
        '200':
          description: Chat response (stream or json)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ChatResponse'
            text/event-stream:
              schema:
                type: string
                description: Server-Sent Events stream

  # ============================
  # Chat History Management
  # ============================
  /core/chat/history/getHistories:
    post:
      summary: Get Chat History List
      tags:
        - Chat History
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                appId:
                  type: string
                offset:
                  type: integer
                  default: 0
                pageSize:
                  type: integer
                  default: 20
                source:
                  type: string
                  enum: [api]
                  description: "Filter by source (e.g. 'api' for API-created chats)"
      responses:
        '200':
          description: List of chat sessions
  
  /core/chat/history/updateHistory:
    post:
      summary: Update Chat Session (Title/Top)
      tags:
        - Chat History
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - appId
                - chatId
              properties:
                appId:
                  type: string
                chatId:
                  type: string
                customTitle:
                  type: string
                  description: Rename the chat session
                top:
                  type: boolean
                  description: Pin or unpin the chat
      responses:
        '200':
          description: Success

  /core/chat/history/delHistory:
    delete:
      summary: Delete a Chat Session
      tags:
        - Chat History
      parameters:
        - name: appId
          in: query
          required: true
          schema:
            type: string
        - name: chatId
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Success

  /core/chat/getPaginationRecords:
    post:
      summary: Get Records within a Chat Session
      tags:
        - Chat History
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - appId
                - chatId
              properties:
                appId:
                  type: string
                chatId:
                  type: string
                offset:
                  type: integer
                pageSize:
                  type: integer
                loadCustomFeedbacks:
                  type: boolean
      responses:
        '200':
          description: List of messages in the session

  # ============================
  # Dataset (Knowledge Base)
  # ============================
  /core/dataset/create:
    post:
      summary: Create Dataset
      tags:
        - Dataset
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DatasetCreateRequest'
      responses:
        '200':
          description: Dataset ID returned

  /core/dataset/list:
    post:
      summary: List Datasets
      tags:
        - Dataset
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                parentId:
                  type: string
                  nullable: true
                  description: Pass null or empty string for root.
      responses:
        '200':
          description: List of datasets

  /core/dataset/detail:
    get:
      summary: Get Dataset Details
      tags:
        - Dataset
      parameters:
        - name: id
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Dataset detail object

  /core/dataset/delete:
    delete:
      summary: Delete Dataset
      tags:
        - Dataset
      parameters:
        - name: id
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Success

  # ============================
  # Dataset Collection (Files/Data)
  # ============================
  /core/dataset/collection/create/text:
    post:
      summary: Create Collection from Text
      tags:
        - Dataset Collection
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - text
                - datasetId
                - name
                - trainingType
              properties:
                text:
                  type: string
                datasetId:
                  type: string
                name:
                  type: string
                trainingType:
                  type: string
                  enum: [chunk, qa]
                chunkSettingMode:
                  type: string
                  enum: [auto, custom]
      responses:
        '200':
          description: Collection creation result

  /core/dataset/collection/create/link:
    post:
      summary: Create Collection from URL
      tags:
        - Dataset Collection
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - link
                - datasetId
                - trainingType
              properties:
                link:
                  type: string
                datasetId:
                  type: string
                trainingType:
                  type: string
                  enum: [chunk, qa]
                metadata:
                  type: object
                  properties:
                    webPageSelector:
                      type: string
                      description: CSS selector to extract content
      responses:
        '200':
          description: Collection creation result

  /core/dataset/data/pushData:
    post:
      summary: Push Data to Collection
      description: Manually insert Q&A or text data into a collection.
      tags:
        - Dataset Data
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - collectionId
                - data
              properties:
                collectionId:
                  type: string
                trainingType:
                  type: string
                  enum: [chunk, qa]
                data:
                  type: array
                  items:
                    type: object
                    required:
                      - q
                    properties:
                      q:
                        type: string
                        description: Main content / Question
                      a:
                        type: string
                        description: Auxiliary content / Answer
                      indexes:
                        type: array
                        items:
                          type: object
                          properties:
                            text:
                              type: string
      responses:
        '200':
          description: Insertion result

  /core/dataset/searchTest:
    post:
      summary: Search Test
      tags:
        - Dataset Data
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - datasetId
                - text
              properties:
                datasetId:
                  type: string
                text:
                  type: string
                limit:
                  type: integer
                  default: 3000
                similarity:
                  type: number
                  default: 0
                searchMode:
                  type: string
                  enum: [embedding, fullTextReRank]
      responses:
        '200':
          description: Search results

  # ============================
  # Share Link Identity
  # ============================
  /shareAuth/init:
    post:
      summary: Share Link Auth Init (User Defined)
      description: This is an interface YOU implement on your server, not a FastGPT API. FastGPT calls this hook.
      tags:
        - Share Auth Hooks
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                token:
                  type: string
      responses:
        '200':
          description: Returns success and uid

  /shareAuth/start:
    post:
      summary: Share Link Auth Start (User Defined)
      description: Hook called before chat starts.
      tags:
        - Share Auth Hooks
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                token:
                  type: string
                question:
                  type: string
      responses:
        '200':
          description: Returns success status

  /shareAuth/finish:
    post:
      summary: Share Link Auth Finish (User Defined)
      description: Hook called after chat finishes (for billing/logging).
      tags:
        - Share Auth Hooks
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                token:
                  type: string
                responseData:
                  type: array
                  items:
                    type: object
      responses:
        '200':
          description: Acknowledge receipt