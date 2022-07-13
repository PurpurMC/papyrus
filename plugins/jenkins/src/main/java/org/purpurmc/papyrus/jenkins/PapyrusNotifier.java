
package org.purpurmc.papyrus.jenkins;

import com.google.common.collect.Lists;
import com.google.common.io.Files;
import edu.umd.cs.findbugs.annotations.NonNull;
import edu.umd.cs.findbugs.annotations.Nullable;
import hudson.Extension;
import hudson.Launcher;
import hudson.model.AbstractBuild;
import hudson.model.AbstractProject;
import hudson.model.BuildListener;
import hudson.scm.ChangeLogSet;
import hudson.tasks.BuildStepDescriptor;
import hudson.tasks.Mailer;
import hudson.tasks.Publisher;
import hudson.tasks.Recorder;
import net.sf.json.JSONObject;
import org.kohsuke.stapler.DataBoundConstructor;
import org.kohsuke.stapler.StaplerRequest;
import org.purpurmc.papyrus.jenkins.uploader.PapyrusUploader;
import org.purpurmc.papyrus.jenkins.uploader.payload.CreateBuildPayload;
import org.purpurmc.papyrus.jenkins.uploader.payload.UploadFilePayload;
import org.purpurmc.papyrus.jenkins.util.Result;

import java.io.File;
import java.io.IOException;
import java.util.List;

public class PapyrusNotifier extends Recorder {
    private final String url;
    private final String key;
    private final String project;
    private final String version;
    private final String file;

    @DataBoundConstructor
    public PapyrusNotifier(String url, String key, String project, String version, String file) {
        /*
        this.url = url;
        this.key = key;
        this.project = project;
        this.version = version;
        this.file = file;
         */

        // todo: this is a hacky workaround of a bug in Jenkins in development
        this.url = "http://127.0.0.1:8000";
        this.key = "key";
        this.project = "purpur";
        this.version = "1.19";
        this.file = "output.txt";
    }

    @Override
    public boolean perform(AbstractBuild<?, ?> build, Launcher launcher, BuildListener listener) throws IOException, InterruptedException {
        PapyrusUploader uploader = new PapyrusUploader(url, key, project, version, file);
        File file = new File(build.getWorkspace().child(this.file).getRemote());
        if (!file.exists()) {
            listener.getLogger().println("File " + file + " does not exist");
            return false;
        }

        listener.getLogger().println("Creating build on papyrus...");

        CreateBuildPayload createBuildPayload = new CreateBuildPayload();
        createBuildPayload.setProject(project);
        createBuildPayload.setVersion(version);
        createBuildPayload.setBuild(String.valueOf(build.getNumber()));
        createBuildPayload.setResult(String.valueOf(build.getResult()));
        createBuildPayload.setDuration(build.getDuration());
        createBuildPayload.setTimestamp(build.getTimeInMillis());

        List<CreateBuildPayload.Commit> commits = Lists.newArrayList();
        for (ChangeLogSet.Entry entry : build.getChangeSet()) {
            CreateBuildPayload.Commit commit = new CreateBuildPayload.Commit();
            commit.setAuthor(entry.getAuthor().getFullName());
            commit.setEmail(entry.getAuthor().getProperty(Mailer.UserProperty.class).getEmailAddress());
            commit.setDescription(entry.getMsg());
            commit.setHash(entry.getCommitId());
            commit.setTimestamp(entry.getTimestamp());
            commits.add(commit);
        }

        createBuildPayload.setCommits(commits);

        Result<String, String> createResult = uploader.create(createBuildPayload);
        if (!createResult.isOk()) {
            listener.getLogger().println("Failed to create build on papyrus: " + createResult.getError());
            return false;
        }

        UploadFilePayload uploadFilePayload = new UploadFilePayload();
        uploadFilePayload.setBuildId(createResult.getValue());
        uploadFilePayload.setFile(file);
        uploadFilePayload.setFileExtension(Files.getFileExtension(file.getName()));

        listener.getLogger().println("Uploading file to papyrus...");
        Result<Object, String> uploadResult = uploader.upload(uploadFilePayload);
        if (!uploadResult.isOk()) {
            listener.getLogger().println("Failed to upload file to papyrus: " + uploadResult.getError());
            return false;
        }

        return true;
    }

    @Extension
    public static class DescriptorImpl extends BuildStepDescriptor<Publisher> {

        @Override
        public String getDisplayName() {
            return "Upload build to Papyrus";
        }

        @Override
        public boolean isApplicable(Class<? extends AbstractProject> jobType) {
            return true;
        }

        @Override
        public Publisher newInstance(@Nullable StaplerRequest req, @NonNull JSONObject formData) throws FormException {
            System.out.println("newInstance");
            System.out.println(formData);
            return new PapyrusNotifier(
                    formData.getString("url"),
                    formData.getString("key"),
                    formData.getString("project"),
                    formData.getString("version"),
                    formData.getString("file")
            );
        }
    }
}
