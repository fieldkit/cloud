// @flow

export type APIErrors = {
  code: string;
  detail: string;
  id: string;
  meta: { [key: string]: string };
  status: number;
}

export type APIBaseUser = {
  username: string;
  name: string;
  bio: string;
  email: string;
}

export type APINewUser = {
  ...$Exact<APIBaseUser>;
  password: string;
  invite_token: string;
}

export type APIUser = {
  ...$Exact<APIBaseUser>;
  id: number;
}

export type APIUsers = {
  users: APIUser[]
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

export type APITwitterInputCreateResponse = {|
  location: string;
|};

export type APIMutableInput = {|
  team_id?: number;
  user_id?: number;
  name: string;
|};

export type APIBaseInput = {|
  id: number;
  expedition_id: number;
  ...APIMutableInput;
|};

export type APINewTwitterInput = {|
  name: string
|};

export type APITwitterInput = {|
  ...APIBaseInput;
  screen_name: string;
  twitter_account_id: number;
|};

export type APIInputs = {
  twitter_accounts?: APITwitterInput[],
  fieldkit_inputs?: APIFieldkitInput[]
};

export type APINewFieldkitInput = {|
  name: string
|};

export type APIFieldkitInput = {|
  ...APIBaseInput;
|};

export type APIFieldkitInputs = {|
  fieldkit_inputs: APIFieldkitInput[];
|};

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

export type APIBaseMember = {|
  role: string;
|}

export type APINewMember = {
  ...$Exact<APIBaseMember>;
  user_id: number;
};

export type APIMember = {
  team_id: number;
   ...$Exact<APINewMember>;
};

export type APIMembers = {
  members: APIMember[]
}
