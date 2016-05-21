using Microsoft.Owin;
using Owin;

[assembly: OwinStartupAttribute(typeof(Guardian.Web.Startup))]
namespace Guardian.Web
{
    public partial class Startup
    {
        public void Configuration(IAppBuilder app)
        {
            ConfigureAuth(app);
        }
    }
}
