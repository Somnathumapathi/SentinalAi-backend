-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create auth.users table (handled by Supabase Auth)
-- Note: This is just for reference, Supabase Auth handles this automatically

-- Create organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    owner_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create github_installations table
CREATE TABLE github_installations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    installation_id BIGINT NOT NULL UNIQUE,
    account_id BIGINT NOT NULL,
    account_type TEXT NOT NULL,
    account_login TEXT NOT NULL,
    repository_selection TEXT NOT NULL,
    access_tokens_url TEXT NOT NULL,
    repositories_url TEXT NOT NULL,
    html_url TEXT NOT NULL,
    app_id BIGINT NOT NULL,
    target_id BIGINT NOT NULL,
    target_type TEXT NOT NULL,
    permissions JSONB NOT NULL,
    events JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create github_repositories table
CREATE TABLE github_repositories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    installation_id UUID REFERENCES github_installations(id) ON DELETE CASCADE,
    repo_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    private BOOLEAN NOT NULL,
    html_url TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    UNIQUE(installation_id, repo_id)
);

-- Create RLS policies
ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE github_installations ENABLE ROW LEVEL SECURITY;
ALTER TABLE github_repositories ENABLE ROW LEVEL SECURITY;

-- Organizations policies
CREATE POLICY "Users can view their own organizations"
    ON organizations FOR SELECT
    TO authenticated
    USING (owner_id = auth.uid());

CREATE POLICY "Users can create organizations"
    ON organizations FOR INSERT
    TO authenticated
    WITH CHECK (owner_id = auth.uid());

CREATE POLICY "Users can update their own organizations"
    ON organizations FOR UPDATE
    TO authenticated
    USING (owner_id = auth.uid());

-- GitHub installations policies
CREATE POLICY "Users can view their organization's installations"
    ON github_installations FOR SELECT
    TO authenticated
    USING (organization_id IN (
        SELECT id FROM organizations WHERE owner_id = auth.uid()
    ));

CREATE POLICY "Users can create installations for their organizations"
    ON github_installations FOR INSERT
    TO authenticated
    WITH CHECK (organization_id IN (
        SELECT id FROM organizations WHERE owner_id = auth.uid()
    ));

-- GitHub repositories policies
CREATE POLICY "Users can view their organization's repositories"
    ON github_repositories FOR SELECT
    TO authenticated
    USING (installation_id IN (
        SELECT id FROM github_installations WHERE organization_id IN (
            SELECT id FROM organizations WHERE owner_id = auth.uid()
        )
    ));

-- Create functions for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = TIMEZONE('utc'::text, NOW());
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updating timestamps
CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_github_installations_updated_at
    BEFORE UPDATE ON github_installations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_github_repositories_updated_at
    BEFORE UPDATE ON github_repositories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 