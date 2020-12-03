FROM gradle:jdk8-openj9
WORKDIR /app
COPY . .

RUN gradle build --no-daemon
RUN cp /app/build/libs/*.jar /app/demo.jar

EXPOSE 8080
CMD ['java','-jar','/app/demo.jar']
