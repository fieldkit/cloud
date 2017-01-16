# fieldkit

## API

### User

```
/api/user/sign-up
```
* invite
* email
* password
* username

```
/api/user/validate
```
* token

```
/api/user/sign-in
```
* username
* password

```
/api/user/current
```

### Project

```
/api/projects
```

```
/api/projects/add
```
* name

```
/api/project/{project}
```

### Expedition

```
/api/project/{project}/expeditions
```

```
/api/project/{project}/expeditions/add
```
* name

```
/api/project/{project}/expedition/{expedition}
```

### Input

```
/api/project/{project}/expedition/{expedition}/inputs
```

```
/api/project/{project}/expedition/{expedition}/inputs/add
```
* name

```
/api/project/{project}/expedition/{expedition}/input/{id}
```

## Frontend

### Visualization App
From `/frontend`

#### install
```
npm install
npm run-script dll
```

#### develop
```
npm start
```
Runs on localhost:8000

#### build
```
npm run-script build
```
Builds app in `/frontend/dist`

### Admin App
From `/admin`

#### install
```
npm install
npm run-script dll
```

#### develop
```
npm start
```
Runs on localhost:8000

#### build
```
npm run-script build
```
Builds app in `/admin/dist`

