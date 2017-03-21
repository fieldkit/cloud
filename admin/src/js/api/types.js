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

// TODO: not officially the API spec, revise when it's there
export type APIUserProfile = {
  ...$Exact<APIUser>;
  name: string;
  bio: string;
  email: string;
  profile_picture: string;
}

// TODO: not officially the API spec, revise when it's there
export type APIPasswordChange = {
  newPassword: string;
};

export type APINewProject = {
  name: string;
  slug: string;
  description: string;
  }

export type APIProject = {
  id: number;
  ...$Exact<APINewProject>;
};

export type APIProjects = {
  projects: APIProject[];
}

export type APINewExpedition = {
  name: string;
  slug: string;
  description: string;
};

export type APIExpedition = {
  id: number;
  ...$Exact<APINewExpedition>;
};

export type APIExpeditions = {
  expeditions: APIExpedition[]
}

export type APIInputType = 'webhook' | 'twitter';

export type APINewInput = {
  name: string;
  type: APIInputType;
  active: boolean;
}

export type APIInput = {
  id: number;
  expedtion_id: number;
  ...$Exact<APINewInput>;
};

export type APIInputs = {
  inputs: APIInput[]
}

export type APINewTeam = {
  name: string;
  slug: string;
  description: string;
};

export type APITeam = {
  id: number;
  ...$Exact<APINewTeam>;
};

export type APITeams = {
  teams: APITeam[]
}

export type APINewAdministrator = {
  user_id: number;
};

export type APIAdministrator = {
  project_id: number;
  ...$Exact<APINewAdministrator>;
};

export type APIAdministrators = {
  administrators: APIAdministrator[]
}

export type APINewMember = {
  user_id: number;
  role: string;
};

export type APIMember = {
  team_id: number;
   ...$Exact<APINewMember>;
};

export type APIMembers = {
  members: APIMember[]
}
