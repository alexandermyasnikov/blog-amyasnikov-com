-- ...

--
-- Name: tasks; Type: TABLE; Schema: taskrepository; Owner: -
--

CREATE TABLE taskrepository.tasks (
    task_id uuid NOT NULL,
    request_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    status text NOT NULL,
    data jsonb NOT NULL
);

--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: taskrepository; Owner: -
--

ALTER TABLE ONLY taskrepository.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (task_id);

--
-- Name: idx_tasks_incomplete; Type: INDEX; Schema: taskrepository; Owner: -
--

CREATE INDEX idx_tasks_incomplete ON taskrepository.tasks USING btree (status) WHERE (status <> 'completed'::text);

-- ...
