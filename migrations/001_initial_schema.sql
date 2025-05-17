-- Create organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    plan TEXT NOT NULL CHECK (plan IN ('free', 'pro', 'enterprise')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'active', 'cancelled', 'expired')),
    razorpay_subscription_id TEXT NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create github_integrations table
CREATE TABLE github_integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    installation_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create aws_integrations table
CREATE TABLE aws_integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    access_key_id TEXT NOT NULL,
    secret_access_key TEXT NOT NULL,
    region TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc'::text, NOW()) NOT NULL
);

-- Create RLS policies
ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE github_integrations ENABLE ROW LEVEL SECURITY;
ALTER TABLE aws_integrations ENABLE ROW LEVEL SECURITY;

-- Create policies for organizations
CREATE POLICY "Organizations are viewable by authenticated users"
    ON organizations FOR SELECT
    TO authenticated
    USING (true);

CREATE POLICY "Organizations can be created by authenticated users"
    ON organizations FOR INSERT
    TO authenticated
    WITH CHECK (true);

CREATE POLICY "Organizations can be updated by their owners"
    ON organizations FOR UPDATE
    TO authenticated
    USING (auth.uid() = id);

-- Create policies for subscriptions
CREATE POLICY "Subscriptions are viewable by organization members"
    ON subscriptions FOR SELECT
    TO authenticated
    USING (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

CREATE POLICY "Subscriptions can be created by organization owners"
    ON subscriptions FOR INSERT
    TO authenticated
    WITH CHECK (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

CREATE POLICY "Subscriptions can be updated by organization owners"
    ON subscriptions FOR UPDATE
    TO authenticated
    USING (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

-- Create policies for github_integrations
CREATE POLICY "GitHub integrations are viewable by organization members"
    ON github_integrations FOR SELECT
    TO authenticated
    USING (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

CREATE POLICY "GitHub integrations can be created by organization owners"
    ON github_integrations FOR INSERT
    TO authenticated
    WITH CHECK (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

-- Create policies for aws_integrations
CREATE POLICY "AWS integrations are viewable by organization members"
    ON aws_integrations FOR SELECT
    TO authenticated
    USING (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
    ));

CREATE POLICY "AWS integrations can be created by organization owners"
    ON aws_integrations FOR INSERT
    TO authenticated
    WITH CHECK (organization_id IN (
        SELECT id FROM organizations WHERE id = organization_id
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

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_github_integrations_updated_at
    BEFORE UPDATE ON github_integrations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_aws_integrations_updated_at
    BEFORE UPDATE ON aws_integrations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 