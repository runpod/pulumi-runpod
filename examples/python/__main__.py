import pulumi
import pulumi_runpod as runpod

my_template = runpod.Template("myTemplate",
    name=f"my-pulumi-template-{pulumi.get_stack()}",
    image_name="runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel-ubuntu22.04",
    container_disk_in_gb=20,
    volume_in_gb=20,
    start_ssh=True,
)

pulumi.export("templateId", my_template.template_id)
