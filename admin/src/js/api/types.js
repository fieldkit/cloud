// @flow

export type APIErrors = {
  code: string;
  detail: string;
  id: string;
  meta: { [key: string]: string };
  status: number;
}

export type APINewUser = {
  email: string;
  username: string;
  password: string;
  invite_token: string;
}

export type APIUser = {
  id: number;
  username: string;
}

export type APIUsers = {
  users: APIUser[]
}

export type APINewProject = {
  name: string;
  slug: string;
  description: string;
}

export type APIProject = {
  id: number;
} & APINewProject;

export type APIProjects = {
  projects: APIProject[];
}

export type APINewExpedition = {
  name: string;
  slug: string;
  description: string;
}

export type APIExpedition = {
  id: number;
} & APINewExpedition;

export type APIExpeditions = {
  expeditions: APIExpedition[]
}

export type APINewInput = {
  name: string;
  slug: string;
  description: string;
}

export type APIInput = {
  id: number;
} & APINewInput;

export type APIInputs = {
  inputs: APIInput[]
}

export type APINewTeam = {
  name: string;
  slug: string;
  description: string;
}

export type APITeam = {
  id: string;
} & APINewTeam;

export type APITeams = {
  teams: APITeam[]
}

export type APINewAdministrator = {
  user_id: number;
}

export type APIAdministrator = {
  project_id: number;
} & APINewAdministrator;

export type APIAdministrators = {
  administrators: APIAdministrator[]
}

export type APINewMember = {
  user_id: number;
  role: string;
};

export type APIMember = {
  team_id: number;
} & APINewMember;

export type APIMembers = {
  members: APIMember[]
}
